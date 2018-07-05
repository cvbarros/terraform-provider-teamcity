<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#macro testInfo tests>
  <#-- @ftlvariable name="tests" type="java.util.List" -->
  <#compress>
    <#if tests?size == 1>
      Test ${tests[0].name} is
    <#else>
      ${tests?size} tests are
    </#if>
  </#compress>
</#macro>

<#macro inScope scopeBean>
  <#-- @ftlvariable name="scopeBean" type="jetbrains.buildServer.notification.TemplateMessageBuilder.MuteScopeBean" -->
  <#compress>
    <#if scopeBean.inProject>
      in project ${scopeBean.project.fullName?html}
    </#if>
    <#if scopeBean.inBuildType>
      in <#list scopeBean.buildTypes as bt>${bt.name?html}<#if bt_has_next>, </#if></#list> (in ${scopeBean.buildTypesProject.fullName?html})
    </#if>
    <#if scopeBean.inBuild>
      in build #${scopeBean.build.buildNumber}
    </#if>
  </#compress>
</#macro>

<#macro unmute unmuteModeBean object>
  <#-- @ftlvariable name="unmuteModeBean" type="jetbrains.buildServer.notification.TemplateMessageBuilder.UnmuteModeBean" -->
  <#compress>
    <#if unmuteModeBean.manually>
      The ${object}s will not be unmuted automatically.
    </#if>
    <#if unmuteModeBean.whenFixed>
      Each ${object} will be unmuted automatically when it is fixed.
    </#if>
    <#if unmuteModeBean.byTime>
      The ${object}s will be unmuted automatically on ${unmuteModeBean.unmuteTime}.
    </#if>
  </#compress>
</#macro>

<#macro unmutedReason unmuteModeBean scopeBean object>
  <#-- @ftlvariable name="unmuteModeBean" type="jetbrains.buildServer.notification.TemplateMessageBuilder.UnmuteModeBean" -->
  <#-- @ftlvariable name="scopeBean" type="jetbrains.buildServer.notification.TemplateMessageBuilder.MuteScopeBean" -->
  <#compress>
    <#if unmuteModeBean.user??>
      Unmute reason: all ${object}s are unmuted manually by ${unmuteModeBean.user.descriptiveName?html}.
    <#else>
      <#if unmuteModeBean.whenFixed>
        Unmute reason: all ${object}s passed successfully <@mute.inScope scopeBean/>.
      </#if>
      <#if unmuteModeBean.byTime>
        Unmute reason: automatically on ${unmuteModeBean.unmuteTime}.
      </#if>
    </#if>
  </#compress>
</#macro>

<#macro comment muteInfo>
  <#-- @ftlvariable name="muteInfo" type="jetbrains.buildServer.serverSide.mute.MuteInfo" -->
  <#compress>
    <#if muteInfo.mutingComment??>
      Comment: ${muteInfo.mutingComment?html}
    </#if>
  </#compress>
</#macro>

<#macro buildProblemInfo buildProblems>
  <#-- @ftlvariable name="buildProblems" type="java.util.List" -->
  <#compress>
    <#if buildProblems?size == 1>
      Build problem is
    <#else>
      ${buildProblems?size} build problems are
    </#if>
  </#compress>
</#macro>

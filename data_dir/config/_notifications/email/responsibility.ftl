<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#macro subject responsibility object>
  <#-- @ftlvariable name="responsibility" type="jetbrains.buildServer.responsibility.ResponsibilityEntry" -->
  <#-- @ftlvariable name="object" type="java.lang.Object" -->
  <#compress>
    <#assign doneByAnotherUser=(responsibility.state != "NONE" &&
                                responsibility.reporterUser?? && responsibility.responsibleUser?? &&
                                responsibility.reporterUser.id != responsibility.responsibleUser.id)/>
    <#assign byWhom>by ${responsibility.reporterUser.descriptiveName?html}</#assign>
    <#if responsibility.state.active>
      <#if responsibility.reporterUser?? && responsibility.responsibleUser?? &&
           responsibility.reporterUser.id == responsibility.responsibleUser.id>
        ${responsibility.responsibleUser.descriptiveName?html} started investigation of ${object}
      <#else>
        ${responsibility.responsibleUser.descriptiveName?html} is assigned for investigation of ${object} ${byWhom}
      </#if>
    </#if>
    <#if responsibility.state.fixed>
      ${responsibility.responsibleUser.descriptiveName?html} fixed ${object} <#if doneByAnotherUser>(done ${byWhom})</#if>
    </#if>
    <#if responsibility.state.givenUp>
      ${responsibility.responsibleUser.descriptiveName?html} gave up fixing ${object} <#if doneByAnotherUser>(done ${byWhom})</#if>
    </#if>
  </#compress>
</#macro>

<#macro comment responsibility>
<#-- @ftlvariable name="responsibility" type="jetbrains.buildServer.responsibility.ResponsibilityEntry" -->
<#if responsibility.comment?length != 0>
Comment: ${responsibility.comment?html}
</#if>
</#macro>

<#macro removeMethod responsibility>
<#-- @ftlvariable name="responsibility" type="jetbrains.buildServer.responsibility.ResponsibilityEntry" -->
Resolve: <#if responsibility.removeMethod.manually>manually<#else>automatically when fixed</#if>
</#macro>

<#macro buildTypeInvestigation buildType successful>
  <#-- @ftlvariable name="buildType" type="jetbrains.buildServer.serverSide.SBuildType" -->
  <#compress>
    <#assign object="Build configuration"/>
    <#assign responsibility=buildType.responsibilityInfo/>
    <#assign state=responsibility.state/>
    <#assign removeManually=responsibility.removeMethod.manually/>

    <#if successful>
      <#if state.active && removeManually>
        ${object} is investigated by ${responsibility.responsibleUser.descriptiveName?html} (can be resolved manually)
      </#if>
      <#if state.fixed>
        ${object} is fixed by ${responsibility.responsibleUser.descriptiveName?html}
      </#if>
    <#else>
      <#if state.active>
        ${object} is investigated by ${responsibility.responsibleUser.descriptiveName?html} <#if removeManually>(can be resolved manually)</#if>
      </#if>
      <#if state.fixed>
        ${object} is marked as fixed by ${responsibility.responsibleUser.descriptiveName?html}
      </#if>
    </#if>
  </#compress>
</#macro>

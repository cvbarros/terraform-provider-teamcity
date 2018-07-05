<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->
<#macro short_build_info build>
  <#-- @ftlvariable name="build" type="jetbrains.buildServer.serverSide.SBuild" -->
  <#if build.branch??>[${build.branch.displayName}] </#if>#${build.buildNumber}
</#macro>

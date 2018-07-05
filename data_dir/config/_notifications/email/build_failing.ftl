<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#-- @ftlvariable name="build" type="jetbrains.buildServer.serverSide.SBuild" -->
<#-- @ftlvariable name="buildType" type="jetbrains.buildServer.serverSide.SBuildType" -->
<#-- @ftlvariable name="project" type="jetbrains.buildServer.serverSide.SProject" -->

<#import "common.ftl" as common>
<#import "responsibility.ftl" as resp>

<#global subject>[<@common.subjMarker/>, FAILING] Build ${project.fullName} / ${buildType.name} <@common.short_build_info build/></#global>

<#global body>Build ${project.fullName} / ${buildType.name} <@common.short_build_info build/> is failing ${var.buildShortStatusDescription}.
<@resp.buildTypeInvestigation buildType false/>
<#if !build.agentLessBuild>Agent: ${agentName}</#if>
Build results: ${link.buildResultsLink}

${var.buildCompilationErrors}${var.buildFailedTestsErrors}${var.buildChanges}
<@common.footer/></#global>

<#global bodyHtml>
  <div>
    <div>
      Build <b>${project.fullName?html} / ${buildType.name?html}</b> <a href='${link.buildResultsLink}'><@common.short_build_info build/></a> is failing
      ${var.buildShortStatusDescription}
    </div>
    <div><@resp.buildTypeInvestigation buildType false/></div>
    <@common.build_agent build/>
    <@common.build_comment build/>
    <br>
    <@common.build_changes var.changesBean/>
    <@common.compilation_errors var.compilationBean/>
    <@common.test_errors var.failedTestsBean/>
    <@common.footerHtml/>
  </div>
</#global>

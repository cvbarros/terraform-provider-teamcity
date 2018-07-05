<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#import "common.ftl" as common>
<#import "responsibility.ftl" as resp>

<#global subject>[<@common.subjMarker/>, STARTED] Build ${project.fullName} / ${buildType.name} <@common.short_build_info build/></#global>

<#global body>Build ${project.fullName} / ${buildType.name} <@common.short_build_info build/> started<#if build.triggeredBy.triggeredByUser> (triggered by ${build.triggeredBy.user.descriptiveName})</#if>.
<@resp.buildTypeInvestigation buildType false/>
<#if !build.agentLessBuild>Agent: ${agentName}</#if>
Build results: ${link.buildResultsLink}

${var.buildChanges}
<@common.footer/></#global>

<#global bodyHtml>
<div>
  <div>
    Build <b>${project.fullName?html} / ${buildType.name?html}</b> <a href='${link.buildResultsLink}'><@common.short_build_info build/></a> started
    <#if build.triggeredBy.triggeredByUser>(triggered by ${build.triggeredBy.user.descriptiveName?html})</#if>
  </div>
  <div><@resp.buildTypeInvestigation buildType false/></div>
  <@common.build_agent build/>
  <@common.build_comment build/>
  <br>
  <@common.build_changes var.changesBean/>
  <@common.footerHtml/>
</div>
</#global>

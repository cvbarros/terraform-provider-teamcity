<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#import "common.ftl" as common>

<#global subject>[<@common.subjMarker/>, LABELING FAILED] ${project.fullName} / ${buildType.name} <@common.short_build_info build/></#global>

<#global body>Labeling failed for root '${vcsRoot.name}'.
Build results: ${link.buildResultsLink}

${var.labelingFailureStacktrace}
<@common.footer/></#global>

<#global bodyHtml>
<div>
  <div>Labeling failed for root <b>'${vcsRoot.name?html}'</b>.</div>
  <div>More details on <a href='${link.buildResultsLink}'>build results page</a></div>
  <br>
  <code style='${common.codeStyle} ${common.stacktraceStyle}'>${var.labelingFailureStacktrace}</code>
  <@common.footerHtml/>
</div>
</#global>

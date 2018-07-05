<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#import "common.ftl" as common>
<#import "responsibility.ftl" as resp>

<#assign subj><@resp.subject responsibility 'build problems'/></#assign>

<#global subject>[<@common.subjMarker/>, INVESTIGATION] ${subj}</#global>

<#global body>${subj} (${project.fullName}).
<@common.build_problem_list buildProblems/>
<@resp.removeMethod responsibility/>
<@resp.comment responsibility/>
${link.allResponsibilitiesLink}
<@common.footer/></#global>

<#global bodyHtml>
<div>
  <div><@resp.subject responsibility 'build problems'/> (${project.fullName?html}):</div>
  <@common.build_problem_list_html buildProblems/>
  <div><@resp.removeMethod responsibility/></div>
  <div><@resp.comment responsibility/></div>
  <br>
  <div>More information on <a href='${link.allResponsibilitiesLink}'>investigations page</a>.</div>
  <@common.footerHtml/>
</div>
</#global>

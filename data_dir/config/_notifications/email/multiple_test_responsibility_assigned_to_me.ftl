<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#import "common.ftl" as common>
<#import "responsibility.ftl" as resp>

<#global subject>[<@common.subjMarker/>, INVESTIGATION] You are assigned for investigation of tests failure</#global>

<#global body>You are assigned for investigation of tests failure (${project.fullName}):
<@common.test_list testNames/>

Assigned by ${responsibility.reporterUser.descriptiveName}
<@resp.removeMethod responsibility/>
<@resp.comment responsibility/>

${link.myResponsibilitiesLink}
<@common.footer/></#global>

<#global bodyHtml>
<div>
  <div>You are assigned for investigation of tests failure (${project.fullName?html}):</div>
  <@common.test_list_html testNames/>
  <div>Assigned by ${responsibility.reporterUser.descriptiveName?html}.</div>
  <div><@resp.removeMethod responsibility/></div>
  <div><@resp.comment responsibility/></div>
  <br>
  <div>More information on <a href='${link.myResponsibilitiesLink}'>investigations page</a>.</div>
  <@common.footerHtml/>
</div>
</#global>

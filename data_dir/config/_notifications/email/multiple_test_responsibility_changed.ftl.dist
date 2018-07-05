<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#import "common.ftl" as common>
<#import "responsibility.ftl" as resp>

<#assign subj><@resp.subject responsibility 'tests'/></#assign>

<#global subject>[<@common.subjMarker/>, INVESTIGATION] ${subj}</#global>

<#global body>${subj}.
<@common.test_list testNames/>
<@resp.removeMethod responsibility/>
<@resp.comment responsibility/>
${link.allResponsibilitiesLink}
<@common.footer/></#global>

<#global bodyHtml>
<div>
  <div><@resp.subject responsibility 'tests failure'/>:</div>
  <@common.test_list_html testNames/>
  <div><@resp.removeMethod responsibility/></div>
  <div><@resp.comment responsibility/></div>
  <br>
  <div>More information on <a href='${link.allResponsibilitiesLink}'>investigations page</a>.</div>
  <@common.footerHtml/>
</div>
</#global>

<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#import "common.ftl" as common>
<#import "responsibility.ftl" as resp>

<#assign subj><@resp.subject responsibility testName/></#assign>

<#global subject>[<@common.subjMarker/>, INVESTIGATION] ${subj}</#global>

<#global body>${subj} (${project.fullName}).
<@resp.removeMethod responsibility/>
<@resp.comment responsibility/>

${link.testLink}
<@common.footer/></#global>

<#global bodyHtml>
<div>
  <div><@resp.subject responsibility '<b>' + testName + '</b>'/> (${project.fullName?html}).</div>
  <div><@resp.removeMethod responsibility/></div>
  <div><@resp.comment responsibility/></div>
  <br>
  <div>More information at <a href='${link.testLink}'>test details page</a>.</div>
  <@common.footerHtml/>
</div>
</#global>

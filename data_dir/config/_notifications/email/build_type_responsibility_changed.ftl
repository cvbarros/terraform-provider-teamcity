<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#import "common.ftl" as common>
<#import "responsibility.ftl" as resp>

<#assign subj><@resp.subject responsibility '${project.fullName} / ${buildType.name}'/></#assign>

<#global subject>[<@common.subjMarker/>, INVESTIGATION] ${subj}</#global>

<#global body>${subj}.
<@resp.removeMethod responsibility/>
<@resp.comment responsibility/>

${link.buildTypeConfigLink}
<@common.footer/></#global>

<#global bodyHtml>
<div>
  <div><@resp.subject responsibility '${project.fullName?html} / ${buildType.name?html}'/>.</div>
  <div><@resp.removeMethod responsibility/></div>
  <div><@resp.comment responsibility/></div>
  <br>
  <div>More information at <a href='${link.buildTypeConfigLink}'>build configuration page</a>.</div>
  <@common.footerHtml/>
</div>
</#global>

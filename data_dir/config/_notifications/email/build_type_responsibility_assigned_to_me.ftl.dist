<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#import "common.ftl" as common>
<#import "responsibility.ftl" as resp>

<#global subject>[<@common.subjMarker/>, INVESTIGATION] You are assigned for investigation of ${project.fullName} / ${buildType.name}</#global>

<#global body>You are assigned for investigation of a build configuration failure.
${project.fullName} / ${buildType.name}, assigned by ${responsibility.reporterUser.descriptiveName}
<@resp.removeMethod responsibility/>
<@resp.comment responsibility/>

${link.buildTypeConfigLink}
<@common.footer/></#global>

<#global bodyHtml>
<div>
  <div>You are assigned for investigation of a build configuration failure.</div>
  <div><b>${project.fullName?html} / ${buildType.name?html}</b>, assigned by ${responsibility.reporterUser.descriptiveName?html}</div>
  <div><@resp.removeMethod responsibility/></div>
  <div><@resp.comment responsibility/></div>
  <br>
  <div>More information at <a href='${link.buildTypeConfigLink}'>build configuration page</a>.</div>
  <@common.footerHtml/>
</div>
</#global>

<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#import "common.ftl" as common>
<#import "responsibility.ftl" as resp>

<#global subject>[<@common.subjMarker/>, INVESTIGATION] You are assigned for investigation of ${testName} failure</#global>

<#global body>You are assigned for investigation of a test failure.
${testName} (${project.fullName}), assigned by ${responsibility.reporterUser.descriptiveName}
<@resp.removeMethod responsibility/>
<@resp.comment responsibility/>

${link.testLink}
<@common.footer/></#global>

<#global bodyHtml>
<div>
  <div>You are assigned for investigation of a test failure.</div>
  <div><b>${testName?html}</b> (${project.fullName?html}), assigned by ${responsibility.reporterUser.descriptiveName?html}.</div>
  <div><@resp.removeMethod responsibility/></div>
  <div><@resp.comment responsibility/></div>
  <br>
  <div>More information at <a href='${link.testLink}'>test details page</a>.</div>
  <@common.footerHtml/>
</div>
</#global>

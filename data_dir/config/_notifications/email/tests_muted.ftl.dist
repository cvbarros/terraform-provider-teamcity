<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#import "common.ftl" as common>
<#import "mute.ftl" as mute>

<#global subject>[<@common.subjMarker/>, MUTE] <@mute.testInfo tests/> muted <@mute.inScope scopeBean/></#global>

<#global body>The following tests are muted <@mute.inScope scopeBean/>:
<#list tests as test>
  ${test.name}
</#list>

User: ${muteInfo.mutingUser.descriptiveName}
<@mute.comment muteInfo/>
<@mute.unmute unmuteModeBean 'test'/>
${link.mutedProblemsLink}
<@common.footer/></#global>

<#global bodyHtml>
  <div>The following tests are muted <@mute.inScope scopeBean/>:</div>
  <ul>
    <#list tests as test>
      <li>${test.name?html}</li>
    </#list>
  </ul>

  <div>
    User: ${muteInfo.mutingUser.descriptiveName?html}
    <br>
    <@mute.comment muteInfo/>
    <br>
    <@mute.unmute unmuteModeBean 'test'/>
  </div>

  <div>More information on <a href='${link.mutedProblemsLink}'>muted problems page</a>.</div>

  <@common.footerHtml/>
</#global>

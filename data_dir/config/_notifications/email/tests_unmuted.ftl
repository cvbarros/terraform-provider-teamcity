<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#import "common.ftl" as common>
<#import "mute.ftl" as mute>

<#global subject>[<@common.subjMarker/>, MUTE] <@mute.testInfo tests/> unmuted</#global>

<#global body>The following tests are unmuted <@mute.inScope scopeBean/>:
<#list tests as test>
  ${test.name}
</#list>

<@mute.unmutedReason unmuteModeBean scopeBean 'test'/>

<@common.footer/></#global>

<#global bodyHtml>
  <div>The following tests are unmuted <@mute.inScope scopeBean/>:</div>
  <ul>
    <#list tests as test>
      <li>${test.name?html}</li>
    </#list>
  </ul>

  <div>
    <@mute.unmutedReason unmuteModeBean scopeBean 'test'/>
  </div>

  <@common.footerHtml/>
</#global>

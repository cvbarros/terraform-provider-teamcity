<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#import "common.ftl" as common>
<#import "mute.ftl" as mute>

<#global subject>[<@common.subjMarker/>, MUTE] <@mute.buildProblemInfo buildProblems/> unmuted</#global>

<#global body>The following build problems are unmuted <@mute.inScope scopeBean/>:
<@common.build_problem_list buildProblems/>

<@mute.unmutedReason unmuteModeBean scopeBean 'build problem'/>

<@common.footer/></#global>

<#global bodyHtml>
  <div>The following build problems are unmuted <@mute.inScope scopeBean/>:</div>
  <@common.build_problem_list_html buildProblems/>

  <div>
    <@mute.unmutedReason unmuteModeBean scopeBean 'build problem'/>
  </div>

  <@common.footerHtml/>
</#global>

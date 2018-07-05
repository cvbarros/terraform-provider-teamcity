<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#import "common.ftl" as common>
<#import "mute.ftl" as mute>

<#global subject>[<@common.subjMarker/>, MUTE] <@mute.buildProblemInfo buildProblems/> muted <@mute.inScope scopeBean/></#global>

<#global body>The following build problems muted <@mute.inScope scopeBean/>:
<@common.build_problem_list buildProblems/>

User: ${muteInfo.mutingUser.descriptiveName}
<@mute.comment muteInfo/>
<@mute.unmute unmuteModeBean 'build problem'/>
${link.mutedProblemsLink}
<@common.footer/></#global>

<#global bodyHtml>
  <div>The following build problems muted <@mute.inScope scopeBean/>:</div>
  <@common.build_problem_list_html buildProblems/>

  <div>
    User: ${muteInfo.mutingUser.descriptiveName?html}
    <br>
    <@mute.comment muteInfo/>
    <br>
    <@mute.unmute unmuteModeBean 'build problem'/>
  </div>

  <div>More information on <a href='${link.mutedProblemsLink}'>muted problems page</a>.</div>

  <@common.footerHtml/>
</#global>

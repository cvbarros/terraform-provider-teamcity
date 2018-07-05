<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#import "common.ftl" as common>
<#import "mute.ftl" as mute>

<#global link>${link.mutedProblemsLink}</#global>
<#global message><@mute.testInfo tests/> muted <@mute.inScope scopeBean/> by ${muteInfo.mutingUser.descriptiveName}.
<@mute.comment muteInfo/></#global>

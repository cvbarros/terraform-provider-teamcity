<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#import "common.ftl" as common>

<#global message>Build is probably hanging.
${project.fullName} / ${buildType.name} <@common.short_build_info build/><#if !build.agentLessBuild>, agent ${agentName}</#if> ${var.buildShortStatusDescription}
Estimated duration: ${var.buildDurationEstimate}, overtime: ${var.buildDurationOvertime}
Last message was received on: ${var.buildLastMessageTimestamp} (${var.buildTimeSinceLastMessage} ago)
${link.buildResultsLink}</#global>

<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#import "common.ftl" as common>

<#global message>Build is failing.
${project.fullName} / ${buildType.name} <@common.short_build_info build/><#if !build.agentLessBuild>, agent ${agentName}</#if> ${var.buildShortStatusDescription}
${link.buildResultsLink}</#global>

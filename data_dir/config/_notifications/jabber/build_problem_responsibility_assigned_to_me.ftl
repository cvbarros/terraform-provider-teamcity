<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#import "common.ftl" as common>
<#import "responsibility.ftl" as resp>

<#global message>You are assigned for build problem investigation in ${project.fullName}:
${buildProblems?first},

assigned by ${responsibility.reporterUser.descriptiveName}

<@resp.removeMethod responsibility/>
<@resp.comment responsibility/>
${link.myResponsibilitiesLink}</#global>
<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#import "common.ftl" as common>
<#import "responsibility.ftl" as resp>

<#global message>Build problem investigation updated in ${project.fullName}:
<@resp.subject responsibility buildProblems?first/>

<@resp.removeMethod responsibility/>
<@resp.comment responsibility/>
${link.allResponsibilitiesLink}</#global>
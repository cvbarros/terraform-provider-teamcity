<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#import "common.ftl" as common>
<#import "responsibility.ftl" as resp>

<#global link>${link.myResponsibilitiesLink}</#global>
<#global message>You are assigned for investigation of tests failure in ${project.fullName}
(${testNames?first} and ${testNames?size - 1} more),
assigned by ${responsibility.reporterUser.descriptiveName}
<@resp.removeMethod responsibility/>
<@resp.comment responsibility/></#global>

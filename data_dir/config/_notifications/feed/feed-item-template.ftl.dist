<#-- TeamCity Defaut Feed Item Template Sample -->
<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#-- TO MAKE THE TEMPLATE EFFECTIVE, RENAME THE FILE TO feed-item-template.ftl -->
<#-- The template is not used until renamed! -->
<#-- The file feed-item-template.ftl.dist will be rewritten on server restart. -->


<#if dataSetType == "FEED">
<#-- Feed Parameters -->
<#if raw_feedType??><#global feedType="${raw_feedType[0]}"><#else>
<#global feedType="atom_1.0">
</#if>

<#assign build_text>
<#if buildStatuses?seq_contains("SUCCESSFUL") && buildStatuses?seq_contains("FAILED")>Builds<#elseif buildStatuses?seq_contains("SUCCESSFUL")>Successful builds<#elseif buildStatuses?seq_contains("FAILED")>Failed builds</#if></#assign>

<#-- Customized feed title -->
<#global feedTitle>
<#if raw_feedTitle??>${raw_feedTitle[0]}<#else>
<#if itemsTypes?seq_contains("BUILDS") && itemsTypes?seq_contains("CHANGES")>${build_text} and changes<#elseif itemsTypes?seq_contains("CHANGES")>Changes<#elseif itemsTypes?seq_contains("BUILDS")>${build_text}</#if> of <#if buildTypes.size() == 1> ${buildTypes.iterator().next().fullName}<#else>${buildTypes.size()} build configurations.</#if>
</#if>
</#global>

<#global feedDescription>
<#if itemsTypes?seq_contains("BUILDS") && itemsTypes?seq_contains("CHANGES")>${build_text} and changes<#elseif itemsTypes?seq_contains("CHANGES")>Changes<#elseif itemsTypes?seq_contains("BUILDS")>${build_text}</#if> of ${buildTypes.toString()} build configuration(s) of TeamCity server at ${globalLinks.root}. (no more then ${itemsCount} items<#if sinceDate??>, only for items since ${sinceDate?date}</#if>)
</#global>

<#global feedAuthor="TeamCity server">


<#global feedLink><#if buildTypes.size() == 1><#-- TODO: refer to project if was specified -->
<#assign buildTypeId=buildTypes.iterator().next().buildTypeId>${globalLinks.getConfigurationHomePage(buildTypeId)}<#else>
${globalLinks.root}</#if></#global>
</#if>


<#if dataSetType == "FEED_ENTRY_BUILD">
<#-- Feed Build Item Template -->
<#global entryTitle>
Build ${project.fullName} / ${buildType.name} <#if build.branch??>[${build.branch.displayName}] </#if>#${build.buildNumber} <#if build.statusDescriptor.successful>was successful<#else>has failed</#if>
</#global>
<#global entryType="html"/> [!-- atom specification notes this should be "html" for HTML content --]

<#-- Customized entry author -->
<#if build.statusDescriptor.successful>${feedEntry.setAuthor("Successful Build")}<#else>${feedEntry.setAuthor("Failed Build")}</#if>

<#global entryDescription>
<#if build.branch??>
Build branch: <strong>${build.branch.displayName}</strong><br/>
</#if>
Status: <strong>${build.statusDescriptor.text}</strong><br/>
Finished on: <strong>${build.finishDate?datetime}</strong><br/>
Changes in the build: <#if build.containingChanges.size() == 0>none<#else>
<a href="${buildLinks.viewChanges}">${build.containingChanges.size()}</a>
by <#list uniqueCommitters as user>
${user}<#if user_has_next>, </#if>
</#list>
</#if><br/>
Agent: <strong>${build.agentName}</strong><br/>
<#if build.shortStatistics.compilationErrorsCount != 0>Compilation errors: <strong>${build.shortStatistics.compilationErrorsCount}</strong>
<br/></#if>
<#if (build.shortStatistics.passedTestCount != 0 ) || (build.shortStatistics.failedTestCount != 0)>
<#assign delimiter="">
Tests: <strong><#if build.shortStatistics.passedTestCount != 0 >${build.shortStatistics.passedTestCount} passed<#assign delimiter=", "></#if>
  <#if build.shortStatistics.failedTestCount != 0 >${delimiter}${build.shortStatistics.failedTestCount} failed</#if>
  <#if build.shortStatistics.ignoredTestCount != 0 >${delimiter}${build.shortStatistics.ignoredTestCount} ignored</#if>
</strong><br/>
<#if build.shortStatistics.failedTestCount != 0>
Failed tests: <#list build.shortStatistics.failedTests as failedTest>
<a href="${buildLinks.getFailedTestResult(failedTest.test.testId)}">
  ${failedTest.test.name.shortName}<#if failedTest.newFailure> (new)</#if></a><#if failedTest_has_next>, </#if>
<#-- Extract that 20 -->
<#if failedTest_index = 20> and ${build.shortStatistics.failedTests.size()-20} more<#break></#if>
</#list><br/>
</#if>
</#if>
<a href="${buildLinks.viewLog}">Build log</a>
</#global>
</#if>


<#if dataSetType == "FEED_ENTRY_CHANGE">
<#-- Feed Change Item Template -->
<#global entryTitle>
Change "<#if change.description != "">${change.description}<#else>No comment</#if>" by ${change.userName} (${change.changeCount} files)
</#global>
<#global entryType="html"/> [!-- atom specification notes this should be "html" for HTML content --]
<#global entryDescription>
Date: <strong>${change.vcsDate?datetime}</strong><br/>
Changed files:<br/>
<#list change.changes as file>
${file.relativeFileName} <i>${file.changeTypeName}</i><br/>
</#list>
</#global>
    </#if>
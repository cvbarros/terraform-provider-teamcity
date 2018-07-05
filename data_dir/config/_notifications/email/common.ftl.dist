<#-- Uses FreeMarker template syntax, template guide can be found at http://freemarker.org/docs/dgui.html -->

<#assign codeStyle>font-family: monospace; font-family: Menlo, Bitstream Vera Sans Mono, Consolas, Courier New, Courier, monospace; font-size: 12px;</#assign>
<#assign failedStyle>color: #e50000;</#assign>
<#assign stacktraceStyle>color: darkred;</#assign>
<#assign separatorStyle>height: 2px; padding: 0; background: #D6D6D6;</#assign>

<#macro subjMarker>TeamCity</#macro>

<#macro plural val><#if (val > 1 || val == 0)>s</#if></#macro>

<#macro build_agent build>
  <#-- @ftlvariable name="build" type="jetbrains.buildServer.serverSide.SBuild" -->
  <#if !build.agentLessBuild>
  <div>Agent: ${build.agentName?html}</div>
  </#if>
</#macro>

<#macro build_comment build>
  <#-- @ftlvariable name="build" type="jetbrains.buildServer.serverSide.SBuild" -->
  <#if build.buildComment??>
    <div>Build comment: <i>${build.buildComment.comment?html}</i> by ${build.buildComment.user.descriptiveName?html}.</div>
  </#if>
</#macro>

<#macro build_changes bean>
  <#-- @ftlvariable name="bean" type="jetbrains.buildServer.notification.impl.ChangesBean" -->
  <#-- @ftlvariable name="webLinks" type="jetbrains.buildServer.serverSide.WebLinks" -->
  <div>
    <#assign modNum=bean.modificationsNumber/>
    <#if (modNum > 0)>
      <div style="${separatorStyle}"></div>
      <br/>
      <div>
        <#assign changesLink><a href='${webLinks.getViewChangesUrl(bean.build)}'>${modNum} change<@plural modNum/></a></#assign>
        Changes included: ${changesLink}<#if bean.changesClipped>,
        only ${bean.modifications?size} are shown</#if>.
      </div>
      <#list bean.modifications as mod>
        <#assign pers><#if mod.personal>(personal build)</#if></#assign>
        <#assign description=mod.description?html/>
        <#if description?length == 0><#assign description='&lt;no comment&gt;'/></#if>
        <div>
          <#assign modLink><a href='${webLinks.getChangeFilesUrl(mod.id, mod.personal)}'>${mod.changes?size} file<@plural mod.changes?size/></a></#assign>
          Change ${mod.displayVersion} ${pers} by ${mod.userName} (${modLink}):
          <i>${description?replace("(\r?\n|\r)", "<br>", "r")?trim}</i>
        </div>
      </#list>
    </#if>
  </div>
</#macro>

<#macro compilation_errors bean>
  <#-- @ftlvariable name="bean" type="jetbrains.buildServer.notification.impl.CompilationErrorsBean" -->
  <#if bean.hasErrorMessages>
    <div>
      <br />
      <div>Compilation errors</div>
      <div style="${separatorStyle}"></div>
      <br/>
      <#list bean.errorMessages as message>
        <code style='${codeStyle} ${stacktraceStyle}'><pre>${message?html}</pre></code>
        <br>
      </#list>
      <#if bean.messagesClipped>
        <code style='${codeStyle}'>&lt;&lt; Error message is clipped &gt;&gt;</code>
      </#if>
    </div>
  </#if>
</#macro>

<#macro test_errors bean>
  <#-- @ftlvariable name="bean" type="jetbrains.buildServer.notification.impl.FailedTestsErrorsBean" -->
  <#if (bean.failedTestCount > 0)>
    <div>
      <br />
      <div>
        Failed tests summary:
        <span style="${failedStyle}">${bean.failedTestCount}
        <#if (bean.newFailedCount > 0)>(${bean.newFailedCount} new)</#if>
        </span>
        <#if bean.summariesClipped>, ${bean.testsForSummary?size} are shown.
          See all on <a href='${link.buildResultsLink}'>build results page</a>.
        </#if>
      </div>
      <table cellspacing="0" cellpadding="5" border="0">
      <#list bean.testsForSummary as test>
        <#assign detailsLink>
          <#if test_index < bean.testDetails?size><a href="#${test.testRunId}" title="Go to stacktrace">details&nbsp;&raquo;</a></#if>
        </#assign>
        <#if test.test.responsibility??>
          <#assign responsibility=test.test.responsibility>
          <#assign investigation>
            <#if responsibility.state.active><i>investigated by ${responsibility.responsibleUser.descriptiveName?html}</i></#if>
            <#if responsibility.state.fixed><i>marked as fixed by ${responsibility.responsibleUser.descriptiveName?html}</i></#if>
            <#if responsibility.state.givenUp><i>given up by ${responsibility.responsibleUser.descriptiveName?html}</i></#if>
          </#assign>
        <#else>
          <#assign investigation=''>
        </#if>
        <tr>
          <td style='padding-left: 10px;'>
            <code style='${codeStyle}'><#if test.newFailure><b>(new) </b></#if>${test.test.name.asString?html}</code>
          </td>
          <td style="font-size: 12px">
            ${investigation}
          </td>
          <td style="font-size: 12px">
            ${detailsLink}
          </td>
        </tr>
      </#list>
      </table>
      <div>
        <br />
        <div>
          <#if (bean.failedTestCount > bean.testDetails?size)>
            Stacktraces (only ${bean.testDetails?size} are shown):
          <#else>
            Stacktraces:
          </#if>
        </div>
        <#list bean.testDetails as details>
          <a name="${details.test.testRunId}"/>
          <code style='${codeStyle}'><#if details.new><b>(new) </b></#if>${details.testName?html}</code>
          <#if details.details??><code style='${codeStyle} ${stacktraceStyle}'><pre>${details.details?html}</pre></code></#if>
          <br>
        </#list>
      </div>
    </div>
  </#if>
</#macro>

<#macro test_list testNames>
  <#-- @ftlvariable name="testNames" type="java.util.List<String>" -->
  <#list testNames as testName>
* ${testName}
  </#list>
</#macro>

<#macro test_list_html testNames>
  <#-- @ftlvariable name="testNames" type="java.util.List<String>" -->
  <ul>
    <#list testNames as testName>
      <li>${testName?html}</li>
    </#list>
  </ul>
</#macro>

<#macro build_problem_list buildProblems>
  <#-- @ftlvariable name="buildProblems" type="java.util.List<String>" -->
  <#list buildProblems as buildProblem>
* ${buildProblem?html?replace("(\r?\n|\r)", "<br>", "r")?trim}
  </#list>
</#macro>

<#macro build_problem_list_html buildProblems>
  <#-- @ftlvariable name="buildProblems" type="java.util.List<String>" -->
  <ul>
    <#list buildProblems as buildProblem>
      <li>${buildProblem?html?replace("(\r?\n|\r)", "<br>", "r")?trim}</li>
    </#list>
  </ul>
</#macro>

<#macro footer>
============================================================================
Configure email notifications: ${link.editNotificationsLink}
</#macro>

<#macro footerHtml>
<div style='color: #666666; font-size:85%'>
  <br/>
  <div style="${separatorStyle}"></div>
  <br/>
  <a href='${link.editNotificationsLink}'>Configure</a> your email notifications on your settings page.
</div>
</#macro>

<#macro short_build_info build><#if build.branch??>[${build.branch.displayName}] </#if>#${build.buildNumber}</#macro>

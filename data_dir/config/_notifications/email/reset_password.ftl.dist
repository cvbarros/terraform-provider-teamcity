<#-- @ftlvariable name="serverUrl" type="java.lang.String" -->
<#-- @ftlvariable name="resetLinks" type="java.util.List<com.intellij.openapi.util.Pair<jetbrains.buildServer.users.SUser, java.lang.String>>" -->
<#-- @ftlvariable name="registerLink" type="java.lang.String" -->
<#-- @ftlvariable name="canRegister" type="java.lang.Boolean" -->

<#global subject>Reset your TeamCity password</#global>

<#global body>
  Hi,

  This email address "${email}" was used when trying to change the password on the TeamCity server ${serverUrl}.

  <#if resetLinks?size == 0>
  The password change has failed because the email address was not found in our database of registered users.
    <#if canRegister>
      Please use the following link if you want to create a new user account: ${registerLink}
    </#if>
  </#if>

  <#if resetLinks?size == 1>
    This email address is associated with the user '${resetLinks[0].first.username}'; use the following link to reset the password: ${resetLinks[0].second}
  <#else>
    This email address is associated with several registered users, use corresponding links to reset their passwords:
    <#list resetLinks as resetLink>
      Reset user '${resetLink.first.username}' <#if (resetLink.first.lastLoginTimestamp)??>(last logged in ${resetLink.first.lastLoginTimestamp?string("dd MMM yy HH:mm")})</#if> password ${resetLink.second}
    </#list>
  </#if>

  If you did not request a password reset from TeamCity, please ignore this message.
</#global>

<#global bodyHtml>
  <div style="padding: 5px 5px; font-size: 14px; line-height: 20px">
      Hi,

      <p>
          This email address "${email}" was used when trying to reset the password on the TeamCity server <a href="${serverUrl}">${serverUrl}</a>.
      </p>

      <#if resetLinks?size == 0>
        <p>
            The password change has failed because the email address was not found in our database of registered users.
        </p>
        <#if canRegister>
            <p>
              Please use the following link if you want to create a new user account:
            </p>
            <p><a href="${registerLink}">Register</a></p>
        </#if>

      <#elseif resetLinks?size == 1>
        <p>
            This email address is associated with the user '${resetLinks[0].first.username}'; use the following link to reset the password: <br/>
        </p>
        <p><a href="${resetLinks[0].second}">Reset Password</a></p>

      <#else>
        <p>
          This email address is associated with several registered users, use corresponding links to reset their passwords:
        </p>
        <p>
          <#list resetLinks as resetLink>
              <a href="${resetLink.second}">Reset Password</a> for user '${resetLink.first.username}' <#if (resetLink.first.lastLoginTimestamp)??>(last logged in ${resetLink.first.lastLoginTimestamp?string("dd MMM yy HH:mm")})</#if><br/>
          </#list>
        </p>
      </#if>

      <p>
          If you did not request a password reset from TeamCity, please ignore this message.
      </p>
  </div>
</#global>

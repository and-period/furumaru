import {
  AdminGetUserCommand,
  AdminGetUserCommandInput,
  AdminGetUserCommandOutput,
  AttributeType,
  AuthFlowType,
  CognitoIdentityProviderClient,
  InitiateAuthCommand,
  InitiateAuthCommandInput,
  InitiateAuthCommandOutput,
} from '@aws-sdk/client-cognito-identity-provider';
import { CognitoUserPoolTriggerEvent } from 'aws-lambda/trigger/cognito-user-pool-trigger';

const client = new CognitoIdentityProviderClient({ region: process.env.AWS_REGION });

const PREVIOUS_COGNITO_USER_POOL_ID = process.env.PREVIOUS_COGNITO_USER_POOL_ID;
const PREVIOUS_COGNITO_CLIENT_ID = process.env.PREVIOUS_COGNITO_CLIENT_ID;

const triggerSources: string[] = ['UserMigration_Authentication', 'UserMigration_ForgotPassword'] as const;

/**
 * Cognitoプールのユーザー情報を別のユーザープールへ移行する
 * @param {Object} event
 * @returns {Object}
 */
export const lambdaHandler = async (event: CognitoUserPoolTriggerEvent): Promise<CognitoUserPoolTriggerEvent> => {
  console.log('received event', JSON.stringify(event));

  if (!triggerSources.includes(event.triggerSource)) {
    throw new Error('bad trigger source');
  }

  if (event.triggerSource === 'UserMigration_Authentication') {
    let auth: InitiateAuthCommandOutput;
    try {
      auth = await initiateAuth(event.userName, event.request.password);
    } catch (err: any) {
      throw new Error(`failed to initiate auth. err=${err.message}`);
    }
    console.log('success to initiate auth', JSON.stringify(auth));
  }

  let user: AdminGetUserCommandOutput;
  try {
    user = await getUser(event.userName);
  } catch (err: any) {
    throw new Error(`failed to get user. err=${err.message}`);
  }
  console.log('success to get user', JSON.stringify(user));

  event.response.userAttributes = toUserAttributes(user);
  event.response.finalUserStatus = 'CONFIRMED';
  event.response.messageAction = 'SUPPRESS';

  console.log('return event', JSON.stringify(event));
  return event;
};

function initiateAuth(username?: string, password?: string): Promise<InitiateAuthCommandOutput> {
  const input: InitiateAuthCommandInput = {
    ClientId: PREVIOUS_COGNITO_CLIENT_ID,
    AuthFlow: AuthFlowType.USER_PASSWORD_AUTH,
    AuthParameters: {
      USERNAME: username || '',
      PASSWORD: password || '',
    },
  };
  return client.send(new InitiateAuthCommand(input));
}

function getUser(username?: string): Promise<AdminGetUserCommandOutput> {
  const input: AdminGetUserCommandInput = {
    UserPoolId: PREVIOUS_COGNITO_USER_POOL_ID,
    Username: username,
  };
  return client.send(new AdminGetUserCommand(input));
}

function toUserAttributes(user: AdminGetUserCommandOutput): { [key: string]: string } {
  const attributes: { [key: string]: string } = {
    username: user.Username || '',
  };
  user.UserAttributes?.forEach((attr: AttributeType): void => {
    switch (attr.Name) {
      case 'sub':
        attributes.sub = attr.Value || '';
        break;
      case 'email':
        attributes.email = attr.Value || '';
        attributes.email_verified = 'true';
        break;
      case 'phone_number':
        attributes.phone_number = attr.Value || '';
        attributes.phone_number_verified = 'true';
        break;
    }
  });
  return attributes;
}

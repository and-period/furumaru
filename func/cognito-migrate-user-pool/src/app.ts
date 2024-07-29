import {
  AdminGetUserCommand,
  AdminGetUserCommandInput,
  AdminGetUserCommandOutput,
  AttributeType,
  AuthFlowType,
  CognitoIdentityProviderClient,
  GetUserCommand,
  GetUserCommandInput,
  GetUserCommandOutput,
  InitiateAuthCommand,
  InitiateAuthCommandInput,
  InitiateAuthCommandOutput,
} from '@aws-sdk/client-cognito-identity-provider';
import { CognitoUserPoolTriggerEvent } from 'aws-lambda/trigger/cognito-user-pool-trigger';

const client = new CognitoIdentityProviderClient({ region: process.env.AWS_REGION });

const PREVIOUS_COGNITO_USER_POOL_ID = process.env.PREVIOUS_COGNITO_USER_POOL_ID;
const PREVIOUS_COGNITO_CLIENT_ID = process.env.PREVIOUS_COGNITO_CLIENT_ID;

/**
 * Cognitoプールのユーザー情報を別のユーザープールへ移行する
 * @param {Object} event
 * @returns {Object}
 */
export const lambdaHandler = async (event: CognitoUserPoolTriggerEvent): Promise<CognitoUserPoolTriggerEvent> => {
  console.log('received event', JSON.stringify(event));

  try {
    if (event.triggerSource === 'UserMigration_Authentication') {
      return await isUserMigrationAuthenticationTriggerEvent(event);
    }
    if (event.triggerSource === 'UserMigration_ForgotPassword') {
      return await isUserMigrationForgotPasswordTriggerEvent(event);
    }
    throw new Error('unknown trigger source');
  } catch (err: any) {
    console.log('received error', err);
    return err.message || `occurred unknown error. err=${err}`;
  }
};

async function isUserMigrationAuthenticationTriggerEvent(
  event: CognitoUserPoolTriggerEvent,
): Promise<CognitoUserPoolTriggerEvent> {
  if (event.triggerSource !== 'UserMigration_Authentication') {
    throw new Error('bad trigger source');
  }

  let auth: InitiateAuthCommandOutput;
  try {
    const input: InitiateAuthCommandInput = {
      ClientId: PREVIOUS_COGNITO_CLIENT_ID,
      AuthFlow: AuthFlowType.USER_PASSWORD_AUTH,
      AuthParameters: {
        USERNAME: event.userName || '',
        PASSWORD: event.request.password || '',
      },
    };
    auth = await client.send(new InitiateAuthCommand(input));
  } catch (err: any) {
    throw new Error(`failed to initiate auth. err=${err.message}`);
  }

  if (!auth) {
    throw new Error('authenticated user is not found');
  }

  let user: GetUserCommandOutput;
  try {
    const input: GetUserCommandInput = {
      AccessToken: auth.AuthenticationResult?.AccessToken || '',
    };
    user = await client.send(new GetUserCommand(input));
  } catch (err: any) {
    throw new Error(`failed to get user. err=${err.message}`);
  }

  if (!user) {
    throw new Error('user is not found');
  }
  console.log('user found', JSON.stringify(user));

  const attributes: { [key: string]: string } = {};
  user.UserAttributes?.forEach((attr: AttributeType): void => {
    switch (attr.Name) {
      case 'email':
        attributes.email = attr.Value || '';
        attributes.email_verified = 'true';
        break;
      case 'sub':
        attributes.sub = attr.Value || '';
        break;
    }
  });

  event.response.userAttributes = attributes;
  event.response.finalUserStatus = 'CONFIRMED';
  event.response.messageAction = 'SUPPRESS';

  console.log('return event', JSON.stringify(event));
  return event;
}

async function isUserMigrationForgotPasswordTriggerEvent(
  event: CognitoUserPoolTriggerEvent,
): Promise<CognitoUserPoolTriggerEvent> {
  if (event.triggerSource !== 'UserMigration_ForgotPassword') {
    throw new Error('bad trigger source');
  }

  let user: AdminGetUserCommandOutput;
  try {
    const input: AdminGetUserCommandInput = {
      UserPoolId: PREVIOUS_COGNITO_USER_POOL_ID,
      Username: event.userName,
    };
    user = await client.send(new AdminGetUserCommand(input));
  } catch (err: any) {
    throw new Error(`failed to get user. err=${err.message}`);
  }

  if (!user) {
    throw new Error('user is not found');
  }
  console.log('user found', JSON.stringify(user));

  const attributes: { [key: string]: string } = {};
  user.UserAttributes?.forEach((attr: AttributeType): void => {
    switch (attr.Name) {
      case 'email':
        attributes.email = attr.Value || '';
        attributes.email_verified = 'true';
        break;
      case 'sub':
        attributes.sub = attr.Value || '';
        break;
    }
  });

  event.response.userAttributes = attributes;
  event.response.finalUserStatus = 'CONFIRMED';
  event.response.messageAction = 'SUPPRESS';

  console.log('return event', JSON.stringify(event));
  return event;
}

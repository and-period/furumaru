/* tslint:disable */
/* eslint-disable */
/**
 * Marche Online
 * マルシェ購入者用API
 *
 * The version of the OpenAPI document: 0.1.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


import * as runtime from '../runtime';
import type {
  AuthGoogleAccountResponse,
  AuthLineAccountResponse,
  AuthResponse,
  CreateAuthUserWithGoogleRequest,
  CreateAuthUserWithLineRequest,
  ErrorResponse,
  ForgotAuthPasswordRequest,
  RefreshAuthTokenRequest,
  ResetAuthPasswordRequest,
  SignInRequest,
  UpdateAuthPasswordRequest,
} from '../models/index';
import {
    AuthGoogleAccountResponseFromJSON,
    AuthGoogleAccountResponseToJSON,
    AuthLineAccountResponseFromJSON,
    AuthLineAccountResponseToJSON,
    AuthResponseFromJSON,
    AuthResponseToJSON,
    CreateAuthUserWithGoogleRequestFromJSON,
    CreateAuthUserWithGoogleRequestToJSON,
    CreateAuthUserWithLineRequestFromJSON,
    CreateAuthUserWithLineRequestToJSON,
    ErrorResponseFromJSON,
    ErrorResponseToJSON,
    ForgotAuthPasswordRequestFromJSON,
    ForgotAuthPasswordRequestToJSON,
    RefreshAuthTokenRequestFromJSON,
    RefreshAuthTokenRequestToJSON,
    ResetAuthPasswordRequestFromJSON,
    ResetAuthPasswordRequestToJSON,
    SignInRequestFromJSON,
    SignInRequestToJSON,
    UpdateAuthPasswordRequestFromJSON,
    UpdateAuthPasswordRequestToJSON,
} from '../models/index';

export interface V1AuthGoogleAccountRequest {
    state: string;
    redirectUri?: string;
}

export interface V1AuthLineAccountRequest {
    state: string;
    redirectUri?: string;
}

export interface V1CreateAuthUserWithGoogleRequest {
    body: CreateAuthUserWithGoogleRequest;
}

export interface V1CreateAuthUserWithLineRequest {
    body: CreateAuthUserWithLineRequest;
}

export interface V1ForgotAuthPasswordRequest {
    body: ForgotAuthPasswordRequest;
}

export interface V1RefreshAuthTokenRequest {
    body: RefreshAuthTokenRequest;
}

export interface V1ResetAuthPasswordRequest {
    body: ResetAuthPasswordRequest;
}

export interface V1SignInRequest {
    body: SignInRequest;
}

export interface V1UpdateUserPasswordRequest {
    body: UpdateAuthPasswordRequest;
}

/**
 * 
 */
export class AuthApi extends runtime.BaseAPI {

    /**
     * Google認証用URLの発行
     */
    async v1AuthGoogleAccountRaw(requestParameters: V1AuthGoogleAccountRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<AuthGoogleAccountResponse>> {
        if (requestParameters['state'] == null) {
            throw new runtime.RequiredError(
                'state',
                'Required parameter "state" was null or undefined when calling v1AuthGoogleAccount().'
            );
        }

        const queryParameters: any = {};

        if (requestParameters['state'] != null) {
            queryParameters['state'] = requestParameters['state'];
        }

        if (requestParameters['redirectUri'] != null) {
            queryParameters['redirectUri'] = requestParameters['redirectUri'];
        }

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/v1/users/me/google`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => AuthGoogleAccountResponseFromJSON(jsonValue));
    }

    /**
     * Google認証用URLの発行
     */
    async v1AuthGoogleAccount(requestParameters: V1AuthGoogleAccountRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<AuthGoogleAccountResponse> {
        const response = await this.v1AuthGoogleAccountRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * LINE認証用URLの発行
     */
    async v1AuthLineAccountRaw(requestParameters: V1AuthLineAccountRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<AuthLineAccountResponse>> {
        if (requestParameters['state'] == null) {
            throw new runtime.RequiredError(
                'state',
                'Required parameter "state" was null or undefined when calling v1AuthLineAccount().'
            );
        }

        const queryParameters: any = {};

        if (requestParameters['state'] != null) {
            queryParameters['state'] = requestParameters['state'];
        }

        if (requestParameters['redirectUri'] != null) {
            queryParameters['redirectUri'] = requestParameters['redirectUri'];
        }

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/v1/users/me/line`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => AuthLineAccountResponseFromJSON(jsonValue));
    }

    /**
     * LINE認証用URLの発行
     */
    async v1AuthLineAccount(requestParameters: V1AuthLineAccountRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<AuthLineAccountResponse> {
        const response = await this.v1AuthLineAccountRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Googleアカウントの連携
     */
    async v1CreateAuthUserWithGoogleRaw(requestParameters: V1CreateAuthUserWithGoogleRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<AuthResponse>> {
        if (requestParameters['body'] == null) {
            throw new runtime.RequiredError(
                'body',
                'Required parameter "body" was null or undefined when calling v1CreateAuthUserWithGoogle().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/v1/users/me/google`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters['body'] as any,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => AuthResponseFromJSON(jsonValue));
    }

    /**
     * Googleアカウントの連携
     */
    async v1CreateAuthUserWithGoogle(requestParameters: V1CreateAuthUserWithGoogleRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<AuthResponse> {
        const response = await this.v1CreateAuthUserWithGoogleRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * LINEアカウントの連携
     */
    async v1CreateAuthUserWithLineRaw(requestParameters: V1CreateAuthUserWithLineRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<AuthResponse>> {
        if (requestParameters['body'] == null) {
            throw new runtime.RequiredError(
                'body',
                'Required parameter "body" was null or undefined when calling v1CreateAuthUserWithLine().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/v1/users/me/line`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters['body'] as any,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => AuthResponseFromJSON(jsonValue));
    }

    /**
     * LINEアカウントの連携
     */
    async v1CreateAuthUserWithLine(requestParameters: V1CreateAuthUserWithLineRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<AuthResponse> {
        const response = await this.v1CreateAuthUserWithLineRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * パスワードリセット
     */
    async v1ForgotAuthPasswordRaw(requestParameters: V1ForgotAuthPasswordRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<void>> {
        if (requestParameters['body'] == null) {
            throw new runtime.RequiredError(
                'body',
                'Required parameter "body" was null or undefined when calling v1ForgotAuthPassword().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/v1/auth/forgot-password`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters['body'] as any,
        }, initOverrides);

        return new runtime.VoidApiResponse(response);
    }

    /**
     * パスワードリセット
     */
    async v1ForgotAuthPassword(requestParameters: V1ForgotAuthPasswordRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<void> {
        await this.v1ForgotAuthPasswordRaw(requestParameters, initOverrides);
    }

    /**
     * トークン検証
     */
    async v1GetAuthRaw(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<AuthResponse>> {
        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("bearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/v1/auth`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => AuthResponseFromJSON(jsonValue));
    }

    /**
     * トークン検証
     */
    async v1GetAuth(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<AuthResponse> {
        const response = await this.v1GetAuthRaw(initOverrides);
        return await response.value();
    }

    /**
     * トークン更新
     */
    async v1RefreshAuthTokenRaw(requestParameters: V1RefreshAuthTokenRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<AuthResponse>> {
        if (requestParameters['body'] == null) {
            throw new runtime.RequiredError(
                'body',
                'Required parameter "body" was null or undefined when calling v1RefreshAuthToken().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("bearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/v1/auth/refresh-token`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters['body'] as any,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => AuthResponseFromJSON(jsonValue));
    }

    /**
     * トークン更新
     */
    async v1RefreshAuthToken(requestParameters: V1RefreshAuthTokenRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<AuthResponse> {
        const response = await this.v1RefreshAuthTokenRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * パスワードリセット - コード検証
     */
    async v1ResetAuthPasswordRaw(requestParameters: V1ResetAuthPasswordRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<void>> {
        if (requestParameters['body'] == null) {
            throw new runtime.RequiredError(
                'body',
                'Required parameter "body" was null or undefined when calling v1ResetAuthPassword().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/v1/auth/forgot-password/verified`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters['body'] as any,
        }, initOverrides);

        return new runtime.VoidApiResponse(response);
    }

    /**
     * パスワードリセット - コード検証
     */
    async v1ResetAuthPassword(requestParameters: V1ResetAuthPasswordRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<void> {
        await this.v1ResetAuthPasswordRaw(requestParameters, initOverrides);
    }

    /**
     * サインイン
     */
    async v1SignInRaw(requestParameters: V1SignInRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<AuthResponse>> {
        if (requestParameters['body'] == null) {
            throw new runtime.RequiredError(
                'body',
                'Required parameter "body" was null or undefined when calling v1SignIn().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("bearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/v1/auth`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters['body'] as any,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => AuthResponseFromJSON(jsonValue));
    }

    /**
     * サインイン
     */
    async v1SignIn(requestParameters: V1SignInRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<AuthResponse> {
        const response = await this.v1SignInRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * サインアウト
     */
    async v1SignOutRaw(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<void>> {
        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("bearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/v1/auth`,
            method: 'DELETE',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.VoidApiResponse(response);
    }

    /**
     * サインアウト
     */
    async v1SignOut(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<void> {
        await this.v1SignOutRaw(initOverrides);
    }

    /**
     * パスワード更新
     */
    async v1UpdateUserPasswordRaw(requestParameters: V1UpdateUserPasswordRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<void>> {
        if (requestParameters['body'] == null) {
            throw new runtime.RequiredError(
                'body',
                'Required parameter "body" was null or undefined when calling v1UpdateUserPassword().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("bearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/v1/auth/password`,
            method: 'PATCH',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters['body'] as any,
        }, initOverrides);

        return new runtime.VoidApiResponse(response);
    }

    /**
     * パスワード更新
     */
    async v1UpdateUserPassword(requestParameters: V1UpdateUserPasswordRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<void> {
        await this.v1UpdateUserPasswordRaw(requestParameters, initOverrides);
    }

}

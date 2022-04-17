<template>
  <div class="text-center">
    <h2>
      &amp;. Marche
    </h2>
    <v-form>
    <v-row justify="center">
      <v-col
        cols="12"
        sm="12"
        md="6"
      >
        <v-text-field
          label="Email"
          placeholder="Email"
          outlined
          required
          type="email"
          dense
        ></v-text-field>
        <v-text-field
          label="Password"
          placeholder="Password"
          outlined
          required
          dense
        ></v-text-field>
      </v-col>
    </v-row>
    <div>
      <nuxt-link to="/">
        パスワードを忘れた場合
      </nuxt-link>
    </div>
    <v-row justify="center">
      <v-col
        cols="12"
        sm="12"
        md="6"
      >
        <v-btn  
          type="submit"
          block
          class="my-5"
        >ログイン</v-btn>
        <div class="px-5 my-3">
          <v-btn  
            block
          >Log In with Google</v-btn>
        </div>
        <div class="px-5 my-3">
        <v-btn  
          block
        >Log In with Facebook</v-btn>
        </div>
        <div class="px-5 my-3">
        <v-btn  
          block
        >Log In with Line</v-btn>
        </div>
      </v-col>
    </v-row>
    </v-form>
    <p>アカウントをお持ちではないですか？
      <nuxt-link to="/">
        登録する
      </nuxt-link>
    </p>
  </div>
</template>

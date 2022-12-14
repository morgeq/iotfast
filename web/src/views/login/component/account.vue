<template>
  <el-form
    ref="loginForm"
    size="large"
    class="login-content-form"
    :model="ruleForm"
    :rules="formRules"
  >
    <el-form-item class="login-animation1" prop="username">
      <el-input
        type="text"
        :placeholder="$t('message.account.accountPlaceholder1')"
        v-model="ruleForm.username"
        clearable
        autocomplete="off"
      >
        <template #prefix>
          <el-icon class="el-input__icon"><ele-User /></el-icon>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item class="login-animation2" prop="password">
      <el-input
        :type="isShowPassword ? 'text' : 'password'"
        :placeholder="$t('message.account.accountPlaceholder2')"
        v-model="ruleForm.password"
        autocomplete="off"
        @keyup.enter="onSignIn"
      >
        <template #prefix>
          <el-icon class="el-input__icon"><ele-Unlock /></el-icon>
        </template>
        <template #suffix>
          <i
            class="iconfont el-input__icon login-content-password"
            :class="isShowPassword ? 'icon-yincangmima' : 'icon-xianshimima'"
            @click="isShowPassword = !isShowPassword"
          >
          </i>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item class="login-animation3" prop="verifyCode">
      <el-col :span="15">
        <el-input
          type="text"
          maxlength="4"
          :placeholder="$t('message.account.accountPlaceholder3')"
          v-model="ruleForm.verifyCode"
          clearable
          autocomplete="off"
          @keyup.enter="onSignIn"
        >
          <template #prefix>
            <el-icon class="el-input__icon"><ele-Position /></el-icon>
          </template>
        </el-input>
      </el-col>
      <el-col :span="1"></el-col>
      <el-col :span="8">
        <div class="login-content-code">
          <img
            class="login-content-code-img"
            @click="getCaptcha"
            width="130"
            height="38"
            :src="captchaSrc"
            style="cursor: pointer"
          />
        </div>
      </el-col>
    </el-form-item>
    <el-form-item class="login-animation4">
      <el-button
        type="primary"
        class="login-content-submit"
        round
        @click="onSignIn"
        :loading="loading.signIn"
      >
        <span>{{ $t("message.account.accountBtnText") }}</span>
      </el-button>
    </el-form-item>
  </el-form>
</template>

<script lang="ts">
import { ElMessage } from "element-plus";
import Cookies from "js-cookie";
import { storeToRefs } from "pinia";
import {
computed, defineComponent, getCurrentInstance, onMounted, reactive, toRefs
} from "vue";
import { useI18n } from "vue-i18n";
import { useRoute, useRouter } from "vue-router";
import { useLoginApi } from "/@/api/login";
import { initBackEndControlRoutes } from "/@/router/backEnd";
import { initFrontEndControlRoutes } from "/@/router/frontEnd";
import { useThemeConfig } from "/@/stores/themeConfig";
import { useUserInfo } from "/@/stores/userInfo";
import { formatAxis } from "/@/utils/formatTime";
import { NextLoading } from "/@/utils/loading";
import { Session } from "/@/utils/storage";

export default defineComponent({
  name: "loginAccount",
  setup() {
    const { t } = useI18n();
    const storesThemeConfig = useThemeConfig();
    const storeUserInfo = useUserInfo();
    const { themeConfig } = storeToRefs(storesThemeConfig);
    const route = useRoute();
    const router = useRouter();
    const { proxy } = <any>getCurrentInstance();
    const state = reactive({
      isShowPassword: false,
      ruleForm: {
        username: "demo",
        password: "123456",
        verifyCode: "",
        verifyKey: "",
      },
      formRules: {
        username: [{ required: true, trigger: "blur", message: "?????????????????????" }],
        password: [{ required: true, trigger: "blur", message: "??????????????????" }],
        verifyCode: [{ required: true, trigger: "blur", message: "?????????????????????" }],
      },
      loading: {
        signIn: false,
      },
      captchaSrc: "",
    });
    onMounted(() => {
      getCaptcha();
    });
    // ????????????
    const currentTime = computed(() => {
      return formatAxis(new Date());
    });
    const getCaptcha = () => {
      useLoginApi()
        .captcha()
        .then((res: any) => {
          state.captchaSrc = res.data.img;
          state.ruleForm.verifyKey = res.data.key;
          //console.log( res.data);
        });
    };
    // ??????
    const onSignIn = () => {
      // ????????????
      proxy.$refs.loginForm
        .validate((valid: boolean) => {
          if (valid) {
            state.loading.signIn = true;
            useLoginApi()
              .signIn(state.ruleForm)
              .then(async (res: any) => {
                //console.log(res.data);
                const userInfos = res.data.userInfo;
                //userInfos.avatar = proxy.getUpFileUrl(userInfos.avatar)
                // ?????? token ??????????????????
                // Session.set('token', Math.random().toString(36).substr(0));
                Session.set("token", res.data.token);
                // ?????????????????????????????????????????????????????????????????????????????????????????? `/src/stores/userInfo.ts` ?????????????????????????????????????????????
                Cookies.set("userName", state.ruleForm.username);

                // ????????????????????????????????????
                Session.set("userInfo", userInfos);
                // ??????????????????
                Session.set("userMenu", res.data.menuList);
                // ??????????????????
                Session.set("permissions", res.data.permissions);
                // 1????????????????????????(?????????????????????vuex)
                //await store.dispatch('userInfos/setUserInfos', userInfos);
                //await store.dispatch('userInfos/setPermissions',res.data.permissions)
                
                await storeUserInfo.setBackUserInfos(userInfos);
                await storeUserInfo.setPermissions(res.data.permissions);

                //console.log("save data");
                if (!themeConfig.value.isRequestRoutes) {
                  // ?????????????????????2????????????????????????
                  await initFrontEndControlRoutes();
                  signInSuccess();
                } else {
                  // ???????????????????????????isRequestRoutes ??? true??????????????????????????????
                  // ????????????????????????????????? router ??????????????????????????? No match found for location with path "/"
                  await initBackEndControlRoutes();
                  // ????????? initBackEndControlRoutes???????????? signInSuccess
                  signInSuccess();
                }
              })
              .catch(() => {
                state.loading.signIn = false;
                getCaptcha();
              });
          }
        })
        .catch(() => {});
    };
    // ????????????????????????
    const signInSuccess = () => {
      //console.log("login successful");
      // ????????????????????????????????????
      let currentTimeInfo = currentTime.value;
      // ??????????????????????????????
      // ????????????????????????????????? router ??????????????????????????? No match found for location with path "/"
      // ??????????????????????????????????????????/???????????????????????????????????????????????????????????????
      if (route.query?.redirect) {
        router.push({
          path: <string>route.query?.redirect,
          query:
            Object.keys(<string>route.query?.params).length > 0
              ? JSON.parse(<string>route.query?.params)
              : "",
        });
      } else {
        router.push("/");
      }
      // ??????????????????
      // ?????? loading
      state.loading.signIn = true;
      const signInText = t("message.signInText");
      ElMessage.success(`${currentTimeInfo}???${signInText}`);
      // ?????? loading???????????????????????????????????????????????????
      NextLoading.start();
    };
    return {
      onSignIn,
      getCaptcha,
      ...toRefs(state),
    };
  },
});
</script>

<style scoped lang="scss">
.login-content-form {
  margin-top: 20px;
  @for $i from 1 through 4 {
    .login-animation#{$i} {
      opacity: 0;
      animation-name: error-num;
      animation-duration: 0.5s;
      animation-fill-mode: forwards;
      animation-delay: calc($i/10) + s;
    }
  }
  .login-content-password {
    display: inline-block;
    width: 20px;
    cursor: pointer;
    &:hover {
      color: #909399;
    }
  }
  .login-content-code {
    display: flex;
    align-items: center;
    justify-content: space-around;
    .login-content-code-img {
      width: 100%;
      height: 40px;
      line-height: 40px;
      background-color: #ffffff;
      border: 1px solid rgb(220, 223, 230);
      cursor: pointer;
      transition: all ease 0.2s;
      border-radius: 4px;
      user-select: none;
      &:hover {
        border-color: #c0c4cc;
        transition: all ease 0.2s;
      }
    }
  }
  .login-content-submit {
    width: 100%;
    letter-spacing: 2px;
    font-weight: 300;
    margin-top: 15px;
  }
}
</style>

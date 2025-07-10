<script lang="ts" setup>
import type { VbenFormSchema } from '@vben/common-ui';
import type { Recordable } from '@vben/types';
import {
  sendCodeApi,
  registerAndActivateApi,
  type RegisterAndActivateParams,
} from '#/api/core/user';
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';

import { AuthenticationRegister, z } from '@vben/common-ui';
import { $t } from '@vben/locales';
import { ElMessage } from 'element-plus';
import { USERNAME_REGEX, EMAIL_REGEX, PASSWORD_REGEX } from '#/config';

defineOptions({ name: 'Register' });

const loading = ref(false);
const registerRef = ref<InstanceType<typeof AuthenticationRegister>>();
const CODE_LENGTH = 6;

const router = useRouter();
const formSchema = computed((): VbenFormSchema[] => {
  return [
    // 用户名，字符长度不超过10个字符，允许使用中文
    {
      component: 'VbenInput',
      componentProps: {
        placeholder: $t('authentication.usernameTip'),
      },
      fieldName: 'username',
      label: $t('authentication.username'),
      rules: z
        .string()
        .min(3, { message: $t('authentication.usernameLengthTip', [3]) })
        .max(10, { message: $t('authentication.usernameLengthErrorTip', [10]) })
        .refine((v) => USERNAME_REGEX.test(v), {
          message: $t('authentication.usernameFormatErrorTip'),
        }),
    },
    {
      component: 'VbenInput',
      componentProps: {
        placeholder: $t('authentication.emailTip'),
      },
      fieldName: 'email',
      label: $t('authentication.email'),
      rules: z
        .string()
        .min(1, { message: $t('authentication.emailTip') })
        .refine((v) => EMAIL_REGEX.test(v), {
          message: $t('authentication.emailValidErrorTip'),
        }),
    },
    {
      component: 'VbenPinInput',
      componentProps: {
        codeLength: CODE_LENGTH,
        createText: (countdown: number) => {
          const text =
            countdown > 0
              ? $t('authentication.sendText', [countdown])
              : $t('authentication.sendCode');
          return text;
        },
        placeholder: $t('authentication.code'),
        handleSendCode,
      },
      fieldName: 'code',
      label: $t('authentication.code'),
      rules: z.string().length(CODE_LENGTH, {
        message: $t('authentication.codeTip', [CODE_LENGTH]),
      }),
    },
    {
      component: 'VbenInputPassword',
      componentProps: {
        passwordStrength: true,
        placeholder: $t('authentication.password'),
      },
      fieldName: 'password',
      label: $t('authentication.password'),
      renderComponentContent() {
        return {
          strengthText: () => $t('authentication.passwordStrength', [8, 16]),
        };
      },
      rules: z
        .string()
        .min(1, { message: $t('authentication.passwordTip') })
        .refine((value) => PASSWORD_REGEX.test(value), {
          message: $t('authentication.passwordFormatErrorTip'),
        }),
    },
    {
      component: 'VbenInputPassword',
      componentProps: {
        placeholder: $t('authentication.confirmPassword'),
      },
      dependencies: {
        rules(values) {
          const { password } = values;
          return z
            .string({ required_error: $t('authentication.passwordTip') })
            .min(1, { message: $t('authentication.passwordTip') })
            .refine((value) => value === password, {
              message: $t('authentication.confirmPasswordTip'),
            });
        },
        triggerFields: ['password'],
      },
      fieldName: 'confirmPassword',
      label: $t('authentication.confirmPassword'),
    },
  ];
});

async function handleSendCode() {
  const formApi = registerRef.value?.getFormApi();
  if (formApi) {
    const email = await formApi.validateField('email');
    if (!email.value) {
      throw new Error($t('authentication.emailTip'));
    }
    await sendCodeApi(email.value);
  }
}

async function handleSubmit(value: Recordable<RegisterAndActivateParams>) {
  // eslint-disable-next-line no-console
  console.log('register submit:', value);
  await registerAndActivateApi(value as unknown as RegisterAndActivateParams);
  // 提示注册成功 mess
  ElMessage.success($t('authentication.registerSuccess'));
  router.push('/auth/login');
}
</script>

<template>
  <AuthenticationRegister
    ref="registerRef"
    :form-schema="formSchema"
    :loading="loading"
    @submit="handleSubmit"
  />
</template>

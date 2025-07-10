<script lang="ts" setup>
import type { VbenFormSchema } from '@vben/common-ui';
import type { Recordable } from '@vben/types';

import { computed, ref } from 'vue';

import { AuthenticationForgetPassword, z } from '@vben/common-ui';
import { $t } from '@vben/locales';
import { ElMessage } from 'element-plus';
import {
  forgotPasswordSendCodeApi,
  forgotPasswordResetApi,
  type ForgotPasswordResetParams,
} from '#/api';
import { PASSWORD_REGEX } from '#/config';
import { useRouter } from 'vue-router';

defineOptions({ name: 'ForgetPassword' });

const loading = ref(false);
const CODE_LENGTH = 6;
const router = useRouter();

const formRef = ref<InstanceType<typeof AuthenticationForgetPassword>>();

const formSchema = computed((): VbenFormSchema[] => {
  return [
    {
      component: 'VbenInput',
      componentProps: {
        placeholder: 'example@example.com',
      },
      fieldName: 'email',
      label: $t('authentication.email'),
      rules: z
        .string()
        .min(1, { message: $t('authentication.emailTip') })
        .email($t('authentication.emailValidErrorTip')),
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
      fieldName: 'newPassword',
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
          const { newPassword } = values;
          return z
            .string({ required_error: $t('authentication.passwordTip') })
            .min(1, { message: $t('authentication.passwordTip') })
            .refine((value) => value === newPassword, {
              message: $t('authentication.confirmPasswordTip'),
            });
        },
        triggerFields: ['newPassword'],
      },
      fieldName: 'confirmPassword',
      label: $t('authentication.confirmPassword'),
    },
  ];
});

async function handleSendCode() {
  const formApi = formRef.value?.getFormApi();
  if (formApi) {
    const email = await formApi.validateField('email');
    if (!email.value) {
      throw new Error($t('authentication.emailTip'));
    }
    await forgotPasswordSendCodeApi(email.value);
  }
}
async function handleSubmit(value: Recordable<any>) {
  // eslint-disable-next-line no-console
  console.log('reset email:', value);
  await forgotPasswordResetApi(value as unknown as ForgotPasswordResetParams);
  // 提示重置密码成功 mess
  ElMessage.success($t('authentication.resetPasswordSuccess'));
  router.push('/auth/login');
}
</script>

<template>
  <AuthenticationForgetPassword
    ref="formRef"
    :form-schema="formSchema"
    :loading="loading"
    @submit="handleSubmit"
  />
</template>

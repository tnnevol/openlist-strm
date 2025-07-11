<script lang="ts" setup>
import { computed } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { ElButton as Button } from 'element-plus';

import { useVbenForm } from '#/adapter/form';
import { $t } from '#/locales';

import { useSchema } from '../data';

const emit = defineEmits(['success']);
const getTitle = computed(() => {
  return '';
});

const [Form, formApi] = useVbenForm({
  layout: 'vertical',
  schema: useSchema(),
  showDefaultActions: false,
});

function resetForm() {
  formApi.resetForm();
}

const [Modal] = useVbenModal({
  async onConfirm() {},
  onOpenChange() {},
});
</script>

<template>
  <Modal :title="getTitle">
    <Form class="mx-4" />
    <template #prepend-footer>
      <div class="flex-auto">
        <Button type="primary" danger @click="resetForm">
          {{ $t('common.reset') }}
        </Button>
      </div>
    </template>
  </Modal>
</template>

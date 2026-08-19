<template>
    <DrawerPro v-model="drawerVisible" :header="$t('setting.title')" @close="handleClose" size="small">
        <el-form ref="formRef" label-position="top" :model="form" :rules="rules" @submit.prevent v-loading="loading">
            <el-form-item :label="$t('setting.title')" prop="panelName">
                <el-input clearable v-model="form.panelName" />
            </el-form-item>

            <el-divider content-position="left">{{ $t('xpack.setting.setting') }}</el-divider>

            <el-form-item :label="$t('xpack.setting.title')">
                <el-input
                    v-model="form.welcomeMessage"
                    type="textarea"
                    :rows="3"
                    maxlength="200"
                    show-word-limit
                    clearable
                />
                <span class="input-help">{{ $t('xpack.setting.titleHelper') }}</span>
            </el-form-item>

            <el-form-item :label="$t('xpack.setting.logoWithText')">
                <div class="branding-upload">
                    <img
                        v-if="form.logo"
                        :src="form.logo"
                        class="branding-logo-preview"
                        :alt="form.panelName || 'Custom logo'"
                    />
                    <div class="flex gap-2">
                        <el-upload
                            :auto-upload="false"
                            :show-file-list="false"
                            accept="image/png,image/jpeg,image/webp"
                            :on-change="onLogoChange"
                        >
                            <el-button>{{ $t('commons.button.upload') }}</el-button>
                        </el-upload>
                        <el-button v-if="form.logo" @click="form.logo = ''">
                            {{ $t('commons.button.delete') }}
                        </el-button>
                    </div>
                </div>
                <span class="input-help">{{ $t('xpack.setting.logoWithTextHelper') }}</span>
            </el-form-item>

            <el-form-item :label="$t('xpack.setting.favicon')">
                <div class="branding-upload">
                    <img
                        v-if="form.websiteIcon"
                        :src="form.websiteIcon"
                        class="branding-icon-preview"
                        alt="Website icon"
                    />
                    <div class="flex gap-2">
                        <el-upload
                            :auto-upload="false"
                            :show-file-list="false"
                            accept="image/png,image/webp,image/x-icon,image/vnd.microsoft.icon,.ico"
                            :on-change="onWebsiteIconChange"
                        >
                            <el-button>{{ $t('commons.button.upload') }}</el-button>
                        </el-upload>
                        <el-button v-if="form.websiteIcon" @click="form.websiteIcon = ''">
                            {{ $t('commons.button.delete') }}
                        </el-button>
                    </div>
                </div>
                <span class="input-help">{{ $t('xpack.setting.faviconHelper') }}</span>
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button @click="drawerVisible = false">{{ $t('commons.button.cancel') }}</el-button>
            <el-button :disabled="loading" type="primary" @click="onSavePanelName(formRef)">
                {{ $t('commons.button.confirm') }}
            </el-button>
        </template>
    </DrawerPro>
</template>
<script lang="ts" setup>
import { reactive, ref } from 'vue';
import i18n from '@/lang';
import { MsgError, MsgSuccess } from '@/utils/message';
import { updateSetting } from '@/api/modules/setting';
import type { FormInstance, UploadFile } from 'element-plus';
import { useGlobalStore } from '@/composables/useGlobalStore';
import { applyWebsiteIcon, getLoginBranding, updateLoginBranding } from '@/utils/branding';
const { themeConfig } = useGlobalStore();

const emit = defineEmits<{ (e: 'search'): void }>();

interface DialogProps {
    panelName: string;
}
const drawerVisible = ref();
const loading = ref(false);

const form = reactive({
    panelName: '',
    welcomeMessage: '',
    logo: '',
    websiteIcon: '',
});
const rules = reactive({
    panelName: [{ validator: checkPanelName, trigger: 'blur', required: true }],
});

function checkPanelName(rule: any, value: any, callback: any) {
    if (value === '') {
        return callback(new Error(i18n.global.t('setting.titleHelper')));
    }
    const reg = /^[a-zA-Z0-9\u4e00-\u9fa5 .,:!@#%&^*_+[\]{}~\-=?，。！｜？：；「」『』【】（）《》·]{3,30}$/;
    if (!reg.test(value)) {
        return callback(new Error(i18n.global.t('setting.titleHelper')));
    }
    callback();
}

const formRef = ref<FormInstance>();

const loadBranding = async () => {
    try {
        const res = await getLoginBranding();
        form.welcomeMessage = res.data.welcomeMessage || '';
        form.logo = res.data.logo || '';
        form.websiteIcon = res.data.websiteIcon || '';
    } catch {
        form.welcomeMessage = '';
        form.logo = '';
        form.websiteIcon = '';
    }
};

const acceptParams = (params: DialogProps): void => {
    form.panelName = params.panelName;
    drawerVisible.value = true;
    loadBranding();
};

type BrandingImageKey = 'logo' | 'websiteIcon';

const readBrandingImage = (key: BrandingImageKey, uploadFile: UploadFile) => {
    const file = uploadFile.raw;
    if (!file) return;

    const isLogo = key === 'logo';
    const isIco = file.name.toLowerCase().endsWith('.ico');
    const allowedLogoTypes = ['image/png', 'image/jpeg', 'image/webp'];
    const allowedIconTypes = ['image/png', 'image/webp', 'image/x-icon', 'image/vnd.microsoft.icon'];
    const isAllowedType = isLogo ? allowedLogoTypes.includes(file.type) : allowedIconTypes.includes(file.type) || isIco;
    if (!isAllowedType) {
        MsgError(i18n.global.t('commons.msg.unSupportType'));
        return;
    }

    const maxSizeMB = isLogo ? 2 : 1;
    if (file.size > maxSizeMB * 1024 * 1024) {
        MsgError(i18n.global.t('commons.msg.unSupportSize', [maxSizeMB]));
        return;
    }

    const reader = new FileReader();
    reader.onload = () => {
        if (typeof reader.result !== 'string') return;
        let result = reader.result;
        if (!isLogo && isIco && !result.startsWith('data:image/')) {
            result = result.replace(/^data:[^;]*;base64,/, 'data:image/x-icon;base64,');
        }
        form[key] = result;
    };
    reader.readAsDataURL(file);
};

const onLogoChange = (file: UploadFile) => readBrandingImage('logo', file);
const onWebsiteIconChange = (file: UploadFile) => readBrandingImage('websiteIcon', file);

const onSavePanelName = async (formEl: FormInstance | undefined) => {
    if (!formEl) return;
    formEl.validate(async (valid) => {
        if (!valid) return;
        loading.value = true;
        try {
            await updateLoginBranding({
                welcomeMessage: form.welcomeMessage,
                logo: form.logo,
                websiteIcon: form.websiteIcon,
            });
            await updateSetting({ key: 'PanelName', value: form.panelName });
            themeConfig.value = {
                ...themeConfig.value,
                panelName: form.panelName,
                title: form.welcomeMessage,
                logo: form.logo,
                logoWithText: form.logo,
                favicon: form.websiteIcon,
            };
            document.title = form.panelName;
            applyWebsiteIcon(form.websiteIcon);
            MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
            drawerVisible.value = false;
            emit('search');
        } finally {
            loading.value = false;
        }
    });
};

const handleClose = () => {
    drawerVisible.value = false;
};

defineExpose({
    acceptParams,
});
</script>

<style scoped lang="scss">
.branding-upload {
    display: flex;
    width: 100%;
    flex-direction: column;
    gap: 8px;
}

.branding-logo-preview {
    max-width: 260px;
    max-height: 64px;
    object-fit: contain;
    align-self: flex-start;
}

.branding-icon-preview {
    width: 32px;
    height: 32px;
    object-fit: contain;
}
</style>

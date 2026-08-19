<template>
    <div class="flex items-center justify-center min-h-screen relative bg-gray-100">
        <div class="absolute inset-0 bg-cover bg-center bg-no-repeat" :style="backgroundStyle"></div>
        <div
            v-if="externalLoginPending"
            v-loading="true"
            class="absolute inset-0 z-20 bg-gray-100"
            aria-busy="true"
        ></div>
        <div
            v-show="!externalLoginPending"
            :style="{ opacity: backgroundOpacity, width: containerWidth, height: containerHeight }"
            class="bg-white shadow-lg relative z-10 border border-gray-200 flex overflow-hidden"
            id="login-container"
        >
            <div class="grid items-stretch w-full" :style="loginGridStyle">
                <div v-if="showLogo" class="flex justify-center" :style="{ height: containerHeight }">
                    <img
                        v-show="imgLoaded"
                        :src="loadImage('loginImage')"
                        class="max-w-full max-h-full object-cover bg-cover bg-center"
                        alt="1panel"
                        @load="onImgLoad"
                        @error="onImgError"
                    />
                </div>
                <div :class="loginFormClass">
                    <div class="w-full h-full flex flex-col">
                        <div
                            v-if="branding.logo || branding.welcomeMessage"
                            class="shrink-0 flex flex-col items-center justify-center gap-1 mb-2 text-center px-4"
                        >
                            <img
                                v-if="branding.logo"
                                :src="branding.logo"
                                class="max-h-10 max-w-[260px] object-contain"
                                :alt="branding.welcomeMessage || themeConfig.panelName || '1Panel'"
                            />
                            <div
                                v-if="branding.welcomeMessage"
                                class="text-sm text-gray-600 leading-5 max-w-full line-clamp-2"
                                :title="branding.welcomeMessage"
                            >
                                {{ branding.welcomeMessage }}
                            </div>
                        </div>
                        <div class="flex-1 min-h-0 overflow-y-auto">
                            <LoginForm ref="loginRef" @external-login-ready="externalLoginPending = false"></LoginForm>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import LoginForm from './components/login-form.vue';
import { ref, onMounted } from 'vue';
import { useGlobalStore } from '@/composables/useGlobalStore';
import { preloadImage } from '@/utils/browser';
import { hasExternalLoginTicket } from '@/utils/external-login';
import { getLoginBranding } from '@/utils/branding';
defineOptions({ name: 'Login' });
const { entrance, isEnterprise, themeConfig } = useGlobalStore();
const backgroundOpacity = ref(1);
const defaultLoginImage = new URL('@/assets/images/1panel-login.jpg', import.meta.url).href;
const defaultEnterpriseLoginImage = new URL('@/assets/images/1panel-login-enterprise.png', import.meta.url).href;
const defaultLoginBgImage = new URL('@/assets/images/1panel-login-bg.jpg', import.meta.url).href;
const loadedLoginImage = ref<string | null>(null);
const loadedBackgroundImage = ref<string | null>(null);
const backgroundStyle = ref<{ backgroundImage?: string; backgroundColor?: string }>({});
const imgLoaded = ref(false);
const currentDefaultLoginImage = computed(() => (isEnterprise.value ? defaultEnterpriseLoginImage : defaultLoginImage));
const branding = ref({
    welcomeMessage: '',
    logo: '',
    websiteIcon: '',
});

const externalLoginPending = ref(hasExternalLoginTicket());

function onImgLoad() {
    imgLoaded.value = true;
}
const mySafetyCode = defineProps({
    code: {
        type: String,
        default: '',
    },
});

const getStatus = async () => {
    let code = mySafetyCode.code;
    if (code != '') {
        entrance.value = code;
    }
};

const loadBranding = async () => {
    try {
        const res = await getLoginBranding();
        branding.value = res.data;
    } catch {
        branding.value = { welcomeMessage: '', logo: '', websiteIcon: '' };
    }
};

const loadImage = (name: string) => {
    const { loginImage, loginBackground, loginBgType } = themeConfig.value;
    if (name === 'loginImage') {
        if (loginImage === 'loginImage') {
            return loadedLoginImage.value || currentDefaultLoginImage.value;
        }
        if (loginImage) {
            return loginImage;
        }
        return currentDefaultLoginImage.value;
    }
    if (name === 'loginBackground') {
        if (loginBgType === 'image') {
            if (loginBackground === 'loginBackground') {
                return loadedBackgroundImage.value || defaultLoginBgImage;
            }
            if (loginBackground) {
                return loginBackground;
            }
            return defaultLoginBgImage;
        }
        if (loginBgType === 'color') {
            return loginBackground;
        }
        return defaultLoginBgImage;
    }
    return '';
};

const onImgError = (event: any) => {
    event.target.src = currentDefaultLoginImage.value;
    imgLoaded.value = true;
};

onMounted(async () => {
    await getStatus();
    await loadBranding();
    const loginImageUrl = `/api/v2/images/loginImage?t=${Date.now()}`;
    const backgroundImageUrl = `/api/v2/images/loginBackground?t=${Date.now()}`;
    if (themeConfig.value.loginImage === 'loginImage') {
        loadedLoginImage.value = await preloadImage(loginImageUrl);
    }
    if (themeConfig.value.loginBgType === 'image' && themeConfig.value.loginBackground === 'loginBackground') {
        loadedBackgroundImage.value = await preloadImage(backgroundImageUrl);
    }
    if (themeConfig.value.loginBgType === 'color') {
        backgroundStyle.value = {
            backgroundColor: themeConfig.value.loginBackground,
        };
    } else {
        const img = new Image();
        const url = loadImage('loginBackground');
        img.onload = () => {
            backgroundStyle.value = {
                backgroundImage: `url(${url})`,
            };
        };
        img.onerror = () => {
            backgroundStyle.value = {
                backgroundImage: `url(${defaultLoginBgImage})`,
            };
        };
        img.src = url;
    }
});

const FIXED_WIDTH = 1000;
const FIXED_HEIGHT = 415;
const useWindowSize = () => {
    const width = ref(window.innerWidth);
    const height = ref(window.innerHeight);

    const updateSize = () => {
        width.value = window.innerWidth;
        height.value = window.innerHeight;
    };

    onMounted(() => window.addEventListener('resize', updateSize));
    onUnmounted(() => window.removeEventListener('resize', updateSize));

    return { width, height };
};
const { width } = useWindowSize();
const showLogo = computed(() => width.value >= FIXED_WIDTH);
const containerWidth = computed(() => `${FIXED_WIDTH}px`);
const containerHeight = computed(() => `${FIXED_HEIGHT}px`);
const loginGridStyle = computed(() => ({
    gridTemplateColumns: showLogo.value ? 'repeat(2, minmax(0, 1fr))' : 'minmax(0, 1fr)',
}));
const loginFormClass = computed(() => {
    return showLogo.value
        ? 'flex items-center justify-center p-4 min-w-0'
        : 'flex items-center justify-center p-4 w-full min-w-0';
});
</script>

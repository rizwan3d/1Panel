import http from '@/api';
import { GlobalStore } from '@/store';

export interface LoginBranding {
    welcomeMessage: string;
    logo: string;
    websiteIcon: string;
}

export const getLoginBranding = () => {
    return http.get<LoginBranding>('/core/auth/branding');
};

export const updateLoginBranding = (params: LoginBranding) => {
    return http.post('/core/settings/branding/update', params);
};

export const applyWebsiteIcon = (websiteIcon: string) => {
    let link = document.querySelector("link[rel*='icon']") as HTMLLinkElement | null;
    if (!link) {
        link = document.createElement('link');
        link.rel = 'shortcut icon';
        document.head.appendChild(link);
    }
    link.href = websiteIcon || '/public/favicon.png';
};

export const loadAndApplyBranding = async () => {
    try {
        const res = await getLoginBranding();
        const globalStore = GlobalStore();
        if (res.data.logo) {
            globalStore.themeConfig.logo = res.data.logo;
            globalStore.themeConfig.logoWithText = res.data.logo;
        }
        if (res.data.websiteIcon) {
            globalStore.themeConfig.favicon = res.data.websiteIcon;
            applyWebsiteIcon(res.data.websiteIcon);
        }
        if (res.data.welcomeMessage) {
            globalStore.themeConfig.title = res.data.welcomeMessage;
        }
        return res.data;
    } catch {
        return null;
    }
};

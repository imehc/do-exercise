import { faker } from '@faker-js/faker';

export interface MenuItem {
  id: number;
  name: string;
  route: string;
  filePath: string;
}

/** 模拟获取管理员菜单 */
export const getAdminMenus = () => {
  return new Promise<MenuItem[]>((resolve) => {
    window.setTimeout(() => {
      const menus: MenuItem[] = Array.from({ length: 3 }, (_, index) => ({
        id: index + 1,
        name: faker.word.sample(),
        route: `/page${index + 1}`,
        filePath: `/page${index + 1}/page.tsx`,
      }));
      resolve(menus);
    }, 1000);
  });
};

/** 模拟获取普通菜单 */
export const getUserMenus = () => {
  return new Promise<MenuItem[]>((resolve) => {
    window.setTimeout(() => {
      const menus: MenuItem[] = Array.from({ length: 1 }, () => ({
        id: 1,
        name: faker.word.sample(),
        route: '/page1',
        filePath: `/page1/page.tsx`,
      }));
      resolve(menus);
    }, 1000);
  });
};
// 用户名，字符长度不超过10个字符，允许使用中文正则
export const USERNAME_REGEX = /^[a-zA-Z0-9\u4e00-\u9fa5]{3,10}$/;

// 邮箱正则
export const EMAIL_REGEX = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

// 密码正则
export const PASSWORD_REGEX =
  /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,16}$/;

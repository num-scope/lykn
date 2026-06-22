declare namespace Api {
  namespace Auth {
    interface LoginResponse {
      access_token: string;
      token_type: string;
      expires_at: string;
      user: UserAccount;
    }

    interface UserAccount {
      id: number;
      username: string;
    }

    interface LoginToken {
      token: string;
      refreshToken: string;
    }

    interface UserInfo {
      userId: string;
      userName: string;
      roles: string[];
      buttons: string[];
    }
  }
}

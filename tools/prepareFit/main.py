import json

from google_auth_oauthlib.flow import InstalledAppFlow


SCOPES = ["https://www.googleapis.com/auth/fitness.activity.read"]


CLIENT_SECRET_FILE = "client_secret.json"


def main() -> None:
    flow = InstalledAppFlow.from_client_secrets_file(
        CLIENT_SECRET_FILE,
        SCOPES,
    )

    creds = flow.run_local_server(
        host="localhost",
        port=0,
        redirect_uri_trailing_slash=False,
    )

    credentials = {
        "type": "authorized_user",
        "client_id": creds.client_id,
        "client_secret": creds.client_secret,
        "refresh_token": creds.refresh_token,
    }

    credentials_json = json.dumps(credentials)
    print(f"GOOGLE_FIT_CREDENTIALS_JSON:\n{credentials_json}")


if __name__ == "__main__":
    main()

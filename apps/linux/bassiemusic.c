#include <gtk/gtk.h>
#include <sys/stat.h>
#include <webkit2/webkit2.h>

static void app_activate(GtkApplication *app) {
    // Create storage folder when it not exists
    char storage_path[255];
    sprintf(storage_path, "%s/.local/share/bassiemusic", getenv("HOME"));
    struct stat sb;
    if (!(stat(storage_path, &sb) == 0 && S_ISDIR(sb.st_mode))) {
        mkdir(storage_path, 0755);
    }

    // Force GTK dark theme
    g_object_set(gtk_settings_get_default(), "gtk-application-prefer-dark-theme", TRUE, NULL);

    // Create window
    GtkWidget *window = gtk_application_window_new(app);
    gtk_window_set_icon_name(GTK_WINDOW(window), "bassiemusic");
    gtk_window_set_default_size(GTK_WINDOW(window), 1280, 720);

    GtkWidget *header_bar = gtk_header_bar_new();
    gtk_header_bar_set_title(GTK_HEADER_BAR(header_bar), "BassieMusic");
    gtk_header_bar_set_show_close_button(GTK_HEADER_BAR(header_bar), TRUE);
    gtk_window_set_titlebar(GTK_WINDOW(window), header_bar);

    // Create webview
    WebKitSettings *settings = webkit_settings_new();
    webkit_settings_set_user_agent(settings, "BassieMusic Linux App v0.1.0");
    GtkWidget *webview = webkit_web_view_new_with_settings(settings);
    gtk_container_add(GTK_CONTAINER(window), GTK_WIDGET(webview));
    WebKitWebsiteDataManager *website_data_manager = webkit_web_view_get_website_data_manager(WEBKIT_WEB_VIEW(webview));
    WebKitCookieManager *cookie_manager = webkit_website_data_manager_get_cookie_manager(website_data_manager);
    char cookies_path[512];
    sprintf(cookies_path, "%s/cookies", storage_path);
    webkit_cookie_manager_set_persistent_storage(cookie_manager, cookies_path, WEBKIT_COOKIE_PERSISTENT_STORAGE_TEXT);
    GdkRGBA background_color = {.red = 0, .blue = 0, .green = 0, .alpha = 0};
    webkit_web_view_set_background_color(WEBKIT_WEB_VIEW(webview), &background_color);
    webkit_web_view_load_uri(WEBKIT_WEB_VIEW(webview), "https://bassiemusic.plaatsoft.nl/");
    gtk_widget_grab_focus(webview);

    // Show window
    gtk_widget_show_all(window);
    gtk_window_present(GTK_WINDOW(window));
}

int main(int argc, char *argv[]) {
    GtkApplication *app = gtk_application_new("nl.plaatsoft.BassieMusic", G_APPLICATION_FLAGS_NONE);
    g_signal_connect(app, "activate", G_CALLBACK(app_activate), NULL);
    return g_application_run(G_APPLICATION(app), argc, argv);
}

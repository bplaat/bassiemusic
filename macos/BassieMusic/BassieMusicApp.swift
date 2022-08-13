//
//  BassieMusicApp.swift
//  BassieMusic
//
//  Created by Bastiaan van der Plaat on 03/07/2022.
//

import SwiftUI

@main
struct BassieMusicApp: App {
    var body: some Scene {
        WindowGroup {
            ContentView().onAppear {
                NSWindow.allowsAutomaticWindowTabbing = false
            }
        }
        .windowToolbarStyle(.unified(showsTitle: false))
        .commands {
            SidebarCommands()
        }
    }
}
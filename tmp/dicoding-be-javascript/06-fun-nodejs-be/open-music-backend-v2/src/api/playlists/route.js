const routes = (handler) => [
    {
        method: 'POST',
        path: '/playlists',
        handler: (request, h) => handler.postPlaylistHandler(request, h),
        options: {
            auth: 'open_music_jwt',
        },
    },
    {
        method: 'GET',
        path: '/playlists',
        handler: (request, h) => handler.getplaylistHandler(request, h),
        options: {
            auth: 'open_music_jwt',
        },
    },
    {
        method: 'DELETE',
        path: '/playlists/{id}',
        handler: (request, h) => handler.deletePlaylistByIdHandler(request, h),
        options: {
            auth: 'open_music_jwt',
        },
    },
    {
        method: 'POST',
        path: '/playlists/{id}/songs',
        handler: (request, h) => handler.postPlaylistSongHandler(request, h),
        options: {
            auth: 'open_music_jwt',
        },
    },
    {
        method: 'GET',
        path: '/playlists/{id}/songs',
        handler: (request, h) => handler.getPlaylistSongHandler(request, h),
        options: {
            auth: 'open_music_jwt',
        },
    },
    {
        method: 'DELETE',
        path: '/playlists/{id}/songs',
        handler: (request, h) => handler.deletePlaylistSongHandler(request, h),
        options: {
            auth: 'open_music_jwt',
        },
    },
    {
        method: 'GET',
        path: '/playlists/{id}/activities',
        handler: (request, h) => handler.getPlaylistActivityHandler(request, h),
        options: {
            auth: 'open_music_jwt',
        },
    },
];

module.exports = routes;

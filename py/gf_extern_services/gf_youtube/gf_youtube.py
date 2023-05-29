# GloFlow application and media management/publishing platform
# Copyright (C) 2023 Ivan Trajkovic
#
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; either version 2 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program; if not, write to the Free Software
# Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA

#---------------------------------------------------------------------------------
def search(p_username_str,
    p_service_client,
    p_results_num_int=20):

    # print("=============================>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", p_username_str)

    #--------------------
    # SEARCH
    r = p_service_client.search().list(
        q=p_username_str,
        part="id,snippet",
        maxResults=p_results_num_int
    ).execute()

    # print(r)

    total_results_for_query_int = r["pageInfo"]["totalResults"]
    items_lst = r["items"]

    #--------------------

    channels_lst = []
    videos_lst   = []

    for item_map in items_lst:

        # print(">>>>>>>>>>>>>>>>>>>", item_map)
        # print("")

        # CHANNEL
        if item_map["id"]["kind"] == "youtube#channel":
            channel_id_str          = item_map["id"]["channelId"]
            snippet_map             = item_map["snippet"]
            channel_title_str       = snippet_map["title"]
            channel_description_str = snippet_map["description"]
            thumbs_map              = snippet_map["thumbnails"]
            thumb_medium_url_str    = thumbs_map["medium"]["url"]

            '''
            print("channel ID", channel_id_str)
            print("title", channel_title_str)
            print("channel snippet", snippet_map)
            print("descr", channel_description_str)
            print("thumb medium", thumb_medium_url_str)
            '''

            channel_map = {
                "id_str":    channel_id_str,
                "title_str": channel_title_str,
                "descr_str": channel_description_str,
                "thumb_medium_url_str": thumb_medium_url_str,
            }
            channels_lst.append(channel_map)

        # VIDEO
        if item_map["id"]["kind"] == "youtube#video":
            video_id_str    = item_map["id"]["videoId"]
            snippet_map     = item_map['snippet']
            channel_id_str  = snippet_map['channelId']
            video_title_str = snippet_map['title']
            description_str = snippet_map["description"]
            thumb_medium_url_str = snippet_map["thumbnails"]["medium"]
            publish_time_str     = snippet_map["publishTime"]

            # print(f"Video ID: {video_id_str}, Video Title: {video_title_str}")

            video_map = {
                "id_str":         video_id_str,
                "channel_id_str": channel_id_str,
                "title_str":      video_title_str,
                "descr_str":      description_str,
                "thumb_medium_url_str": thumb_medium_url_str,
                "publish_time_str":     publish_time_str,
            }

            videos_lst.append(video_map)


    return channels_lst, videos_lst

#---------------------------------------------------------------------------------
def channel_info(p_username_str):

    # IMPORTANT!! - most usernames dont return any results for channels.
    #               in Youtube there is a privacy option for channel owners to disable
    #               their channel info being returned by Youtube API.
    r = p_service_client.channels().list(
            part="id,snippet,contentDetails,statistics",
            forUsername=p_username_str
        ).execute()
    print(r)
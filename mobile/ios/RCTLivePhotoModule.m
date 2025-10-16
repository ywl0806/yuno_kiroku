// RCTLivePhotoModule.m
#import "RCTLivePhotoModule.h"
#import <React/RCTLog.h>
#import <Photos/Photos.h>

@implementation RCTLivePhotoModule

RCT_EXPORT_MODULE(LivePhotoModule);

RCT_EXPORT_METHOD(getLivePhotoMov:(NSString *)phurl resolver:(RCTPromiseResolveBlock)resolve rejecter:(RCTPromiseRejectBlock)reject)
{
    
    // PHAsset 가져오기
    PHFetchResult<PHAsset *> *assets = [PHAsset fetchAssetsWithLocalIdentifiers:@[phurl] options:nil];
    PHAsset *livePhotoAsset = [assets firstObject];
    
  // 라이브 포토의 asset으로 부터 video의 URL을 가져오기
   [self videoUrlForLivePhotoAsset:livePhotoAsset withCompletionBlock:^(NSURL *url) {
       if (url) {
           resolve([url absoluteString]);
       } else {
           reject(@"error", @"Failed to get video URL for live photo", nil);
       }
   }];
}

- (void)videoUrlForLivePhotoAsset:(PHAsset*)asset withCompletionBlock:(void (^)(NSURL* url))completionBlock {
    if (![asset isKindOfClass:[PHAsset class]]) {
        completionBlock(nil);
        return;
    }

    PHAssetResource *videoResource = nil;

    NSArray *assetResources = [PHAssetResource assetResourcesForAsset:asset];
    for (PHAssetResource *resource in assetResources) {
        if (resource.type == PHAssetResourceTypePairedVideo) {
            videoResource = resource;
            break;
        }
    }

    if (!videoResource) {
        completionBlock(nil);
        return;
    }

    NSString *fileName = [NSString stringWithFormat:@"%@.mov", [[NSUUID UUID] UUIDString]];
    NSString *filePath = [NSTemporaryDirectory() stringByAppendingPathComponent:fileName];
    NSURL *fileUrl = [NSURL fileURLWithPath:filePath];

    [[PHAssetResourceManager defaultManager] writeDataForAssetResource:videoResource toFile:fileUrl options:nil completionHandler:^(NSError * _Nullable error) {
        if (error) {
            NSLog(@"Error saving video: %@", error);
            completionBlock(nil);
        } else {
            completionBlock(fileUrl);
        }
    }];
}


@end

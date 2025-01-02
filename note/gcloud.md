# gcloud win-cmd:^ linux cmd:\

# login and set projectID:
gcloud auth login
gcloud config set project deep-rigging-424917-a7

# add role
gcloud projects add-iam-policy-binding deep-rigging-424917-a7 \
    --member="serviceAccount:760527922880-compute@developer.gserviceaccount.com" \
    --role="roles/run.admin"

gcloud projects add-iam-policy-binding deep-rigging-424917-a7 \
    --member="serviceAccount:760527922880-compute@developer.gserviceaccount.com" \
    --role="roles/iam.serviceAccountUser"

gcloud projects add-iam-policy-binding deep-rigging-424917-a7 \
    --member="serviceAccount:760527922880-compute@developer.gserviceaccount.com" \
    --role="roles/cloudbuild.builds.editor"

gcloud projects add-iam-policy-binding deep-rigging-424917-a7 \
    --member="serviceAccount:760527922880-compute@developer.gserviceaccount.com" \
    --role="roles/storage.objectAdmin"

gcloud projects add-iam-policy-binding deep-rigging-424917-a7 \
    --member="serviceAccount:760527922880-compute@developer.gserviceaccount.com" \
    --role="roles/artifactregistry.admin"


# set workload-identity-pools providers
gcloud iam workload-identity-pools providers describe github-actions-provider \
    --workload-identity-pool=github-actions-pool \
    --project=deep-rigging-424917-a7

    
# 添加仓库到 IAM 绑定
gcloud iam service-accounts add-iam-policy-binding 760527922880-compute@developer.gserviceaccount.com \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/projects/760527922880/locations/global/workloadIdentityPools/github-actions-pool/attribute.repository/Rollbytes/mai_go"



# 验证绑定
gcloud iam service-accounts get-iam-policy 760527922880-compute@developer.gserviceaccount.com




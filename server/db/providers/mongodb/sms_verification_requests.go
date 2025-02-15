package mongodb

import (
	"context"
	"time"

	"github.com/authorizerdev/authorizer/server/db/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SMS verification Request
func (p *provider) UpsertSMSRequest(ctx context.Context, smsRequest *models.SMSVerificationRequest) (*models.SMSVerificationRequest, error) {
	smsVerificationRequest, _ := p.GetCodeByPhone(ctx, smsRequest.PhoneNumber)
	shouldCreate := false

	if smsVerificationRequest == nil {
		id := uuid.NewString()

		smsVerificationRequest = &models.SMSVerificationRequest{
			ID: 		id,
			CreatedAt:	time.Now().Unix(),
			Code: smsRequest.Code,
			PhoneNumber: smsRequest.PhoneNumber,
			CodeExpiresAt: smsRequest.CodeExpiresAt,
		}
		shouldCreate = true
	}
	
	smsVerificationRequest.UpdatedAt = time.Now().Unix()
	smsRequestCollection := p.db.Collection(models.Collections.SMSVerificationRequest, options.Collection())

	var err error
	if shouldCreate {
		_, err = smsRequestCollection.InsertOne(ctx, smsVerificationRequest)
	} else {
		_, err = smsRequestCollection.UpdateOne(ctx, bson.M{"phone_number": bson.M{"$eq": smsRequest.PhoneNumber}}, bson.M{"$set": smsVerificationRequest}, options.MergeUpdateOptions())
	}

	if err != nil {
		return nil, err
	}
	
	return smsVerificationRequest, nil
}

func (p *provider) GetCodeByPhone(ctx context.Context, phoneNumber string) (*models.SMSVerificationRequest, error) {
	var smsVerificationRequest models.SMSVerificationRequest

	smsRequestCollection := p.db.Collection(models.Collections.SMSVerificationRequest, options.Collection())
	err := smsRequestCollection.FindOne(ctx, bson.M{"phone_number": phoneNumber}).Decode(&smsVerificationRequest)

	if err != nil {
		return nil, err
	}

	return &smsVerificationRequest, nil
}

func (p *provider) DeleteSMSRequest(ctx context.Context, smsRequest *models.SMSVerificationRequest) error {
	smsVerificationRequests := p.db.Collection(models.Collections.SMSVerificationRequest, options.Collection())
	_, err := smsVerificationRequests.DeleteOne(nil, bson.M{"_id": smsRequest.ID}, options.Delete())
	if err != nil {
		return err
	}

	return nil
}

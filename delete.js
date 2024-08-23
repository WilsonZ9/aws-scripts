process.env.AWS_ACCESS_KEY_ID = ""
process.env.AWS_SECRET_ACCESS_KEY = ""
process.env.AWS_SESSION_TOKEN = ""

const AWS = require('aws-sdk');

AWS.config.update({ region: 'us-east-1' });

const dynamodb = new AWS.DynamoDB.DocumentClient();

const tableName = 'staging_campaigns_practice_enrollment';

const params = {
    TableName: tableName,
};

let count = 0;
dynamodb.scan(params, (err, data) => {
    if (err) {
        console.error('Error scanning table:', err);
    } else {
        const items = data.Items;
        for (const item of items) {
            if (item.campaignEnrollments) {
                for (let i = 0; i < item.campaignEnrollments.length; i++) {
                    const enrollment = item.campaignEnrollments[i];
                    if (enrollment.campaignConfig && enrollment.campaignConfig.appointmentTypeGlobalSetting) {
                        const updateParams = {
                            TableName: tableName,
                            Key: {
                                practiceRowKey: item.practiceRowKey
                            },
                            UpdateExpression: `REMOVE campaignEnrollments[${i}].campaignConfig.appointmentTypeGlobalSetting`
                        };

                        dynamodb.update(updateParams, (updateErr, updateData) => {
                            if (updateErr) {
                                console.error('Error updating item:', updateErr);
                            } else {
                                count++;
                                console.log(`Item updated successfully:${count}`, updateData);
                            }
                        });
                    }
                }
            }
        }
    }
});

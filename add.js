process.env.AWS_ACCESS_KEY_ID = ""
process.env.AWS_SECRET_ACCESS_KEY = ""
process.env.AWS_SESSION_TOKEN = ""

const AWS = require('aws-sdk');

AWS.config.update({ region: 'us-east-1' });

const dynamodb = new AWS.DynamoDB.DocumentClient();

const tableName = 'prod_campaigns_practice_enrollment';

const params = {
    TableName: tableName,
};

let count = 0;
const scanTable = async (params) => {
    dynamodb.scan(params, (err, data) => {
        if (err) {
            console.error('Error scanning table:', err);
        } else {
            const items = data.Items;
            for (const item of items) {
                if (item.campaignEnrollments) {
                    for (let i = 0; i < item.campaignEnrollments.length; i++) {
                        const enrollment = item.campaignEnrollments[i];
                        if (enrollment.campaignConfig && enrollment.campaignType === "appointment-reminders") {
                            let newAppointmentTypeGlobalSettings;
                            if (!enrollment.campaignConfig.appointmentTypeGlobalSettings) {
                                newAppointmentTypeGlobalSettings = {
                                    bloodwork: {
                                        costType: "range",
                                        minCost: 80,
                                        maxCost: 200
                                    }
                                };
                            } else {
                                newAppointmentTypeGlobalSettings = {
                                    bloodwork: {
                                        costType: enrollment.campaignConfig.appointmentTypeGlobalSettings.bloodwork.costType,
                                        minCost: enrollment.campaignConfig.appointmentTypeGlobalSettings.bloodwork.minCost ? Number(enrollment.campaignConfig.appointmentTypeGlobalSettings.bloodwork.minCost) : undefined,
                                        maxCost: enrollment.campaignConfig.appointmentTypeGlobalSettings.bloodwork.maxCost ? Number(enrollment.campaignConfig.appointmentTypeGlobalSettings.bloodwork.maxCost) : undefined,
                                        singleCost: enrollment.campaignConfig.appointmentTypeGlobalSettings.bloodwork.singleCost ? Number(enrollment.campaignConfig.appointmentTypeGlobalSettings.bloodwork.singleCost) : undefined,
                                    }
                                };
                            }
                            enrollment.campaignConfig.appointmentTypeGlobalSettings = newAppointmentTypeGlobalSettings
                            const updateParams = {
                                TableName: tableName,
                                Key: {
                                    practiceRowKey: item.practiceRowKey
                                },
                                UpdateExpression: `set campaignEnrollments[${i}].campaignConfig.appointmentTypeGlobalSettings = :newSetting`,
                                ExpressionAttributeValues: {
                                    ':newSetting': enrollment.campaignConfig.appointmentTypeGlobalSettings
                                }
                            };

                            dynamodb.update(updateParams, (updateErr, updateData) => {
                                if (updateErr) {
                                    console.error('Error updating item:', item.practiceRowKey);
                                } else {
                                    count++;
                                    console.log(`Item updated successfully:${count}`, item.practiceRowKey);
                                }
                            });

                        }
                    }
                }
            }

            // Check if there are more items to scan
            if (data.LastEvaluatedKey) {
                params.ExclusiveStartKey = data.LastEvaluatedKey;
                scanTable(params); // Continue scanning with the new start key
            }
        }
    });
};

scanTable(params);
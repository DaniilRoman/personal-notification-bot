import logging
from typing import Optional

import boto3


class GenericObject(object):
    def __repr__(self):
        return str(self.__dict__)


class StoredItem(GenericObject):
    table_name = "items"
    app_name = "person-notification-bot"

    def __init__(self, app_name: str, item_name: str, item_value: str):
        self.app_name = app_name
        self.item_name = item_name
        self.item_value = item_value


class DynamodbConfig:
    def __init__(self, aws_access_key_id='', aws_secret_access_key='', region_name='', endpoint_url=None):
        self.aws_access_key_id = aws_access_key_id
        self.aws_secret_access_key = aws_secret_access_key
        self.region_name = region_name
        self.endpoint_url = endpoint_url


class ItemStoreService:
    def __init__(self, dynamo_config: DynamodbConfig):
        self.dynamodb_client = boto3.client(
            'dynamodb',
            aws_access_key_id=dynamo_config.aws_access_key_id,
            aws_secret_access_key=dynamo_config.aws_secret_access_key,
            region_name=dynamo_config.region_name,
            endpoint_url=dynamo_config.endpoint_url
        )

        self.dynamodb = boto3.resource(
            'dynamodb',
            aws_access_key_id=dynamo_config.aws_access_key_id,
            aws_secret_access_key=dynamo_config.aws_secret_access_key,
            region_name=dynamo_config.region_name,
            endpoint_url=dynamo_config.endpoint_url
        )

        self.serializer = boto3.dynamodb.types.TypeSerializer()
        self.deserializer = boto3.dynamodb.types.TypeDeserializer()

    def create_item_table(self):
        try:
            self.dynamodb_client.create_table(
                TableName=StoredItem.table_name,
                KeySchema=[
                    {
                        'AttributeName': 'app_name',
                        'KeyType': 'HASH'  # Partition key
                    },
                    {
                        'AttributeName': 'item_name',
                        'KeyType': 'RANGE'  # Sort key
                    }
                ],
                AttributeDefinitions=[
                    {
                        'AttributeName': 'app_name',
                        'AttributeType': 'S'
                    },
                    {
                        'AttributeName': 'item_name',
                        'AttributeType': 'S'
                    },
                ],
                # ProvisionedThroughput={
                #     'ReadCapacityUnits': 5,
                #     'WriteCapacityUnits': 5
                # },
                BillingMode='PAY_PER_REQUEST'
            )
        except self.dynamodb_client.exceptions.ResourceInUseException:
            logging.info(f"Table {StoredItem.table_name} already exists")

    def get_item(self, item_name: str) -> Optional[str]:
        return self.__get_item(StoredItem.app_name, item_name)

    def __get_item(self, app_name: str, item_name: str) -> Optional[str]:
        list_values = self.__get_list_with_error_handling(
            lambda: self.dynamodb_client.execute_statement(
                Statement=f'SELECT * FROM {StoredItem.table_name} WHERE app_name = \'{app_name}\' AND item_name = \'{item_name}\''),
            StoredItem
        )
        if list_values:
            return list_values[0].item_value
        else:
            return None

    def save_item(self, name, value):
        self.__save_item(StoredItem(StoredItem.app_name, name, value))

    ########################################################################################################################

    def __to_dict(self, data: object):
        return {k: self.serializer.serialize(v) for k, v in data.__dict__.items()}

    def __from_dict(self, data: dict):
        return {k: self.deserializer.deserialize(v) for k, v in data.items()}

    def __save_item(self, item):
        return self.__save_with_error_handling(self.dynamodb_client, StoredItem.table_name, item)

    def __save_with_error_handling(self, client, table_name, item):
        return self.__exec_with_error_handling(lambda: client.put_item(
            TableName=table_name,
            Item=self.__to_dict(item)
        ))

    def __exec_with_error_handling(self, specific_operation_func):
        res = specific_operation_func()
        if res["ResponseMetadata"]["HTTPStatusCode"] != 200:
            raise RuntimeError("Cannot delete item")
        return res

    def __get_list_with_error_handling(self, get_func, class_type):
        res = self.__exec_with_error_handling(get_func)
        return [class_type(**self.__from_dict(i)) for i in res["Items"]]

#!/usr/bin/python3

import sys

from typing import List

class CommitMsgValidator:
    MESSAGE_PREFIXES: List[str] = [
        'feat',
        'fix',
        'chore'
    ]

    MESSAGE_PREFIX_SEPARATOR: str = ':'

    def _handle_incorrect_format(self) -> None:
        print('Incorrect format! Message should be formatted as follows:')
        print('<type>: <subject>')

    def _handle_incorrect_prefix(self, prefix: str) -> None:
        print(f'Message prefix "{prefix}" is incorrect!\n')
        print(f'Possible prefixes are: {self.MESSAGE_PREFIXES[0]}, {self.MESSAGE_PREFIXES[1]}, {self.MESSAGE_PREFIXES[2]}')

    def _handle_incorrect_lettercase(self) -> None:
        print('Message should be all lowercase!')

    def _get_message(self) -> str:
        # Git saves commit message in a file, and passes
        # it's name as script's argument
        message_filename = sys.argv[1]
        handler = open(message_filename)

        return handler.read()

    def _is_prefix_correct(self, prefix: str) -> bool:
        return prefix in self.MESSAGE_PREFIXES

    def validate(self) -> None:
        message = self._get_message()

        parts = message.split(self.MESSAGE_PREFIX_SEPARATOR)

        if len(parts) != 2:
            self._handle_incorrect_format()

            return False

        prefix = parts[0]
        subject = parts[1]

        if not self._is_prefix_correct(prefix):
            self._handle_incorrect_prefix(prefix)

            return False
        
        if not subject.islower():
            self._handle_incorrect_lettercase()

            return False

        return True

FAILURE_EXIT_CODE = 1
SUCCESS_EXIT_CODE = 0

validator = CommitMsgValidator()
is_correct = validator.validate()

sys.exit(SUCCESS_EXIT_CODE if is_correct else FAILURE_EXIT_CODE)

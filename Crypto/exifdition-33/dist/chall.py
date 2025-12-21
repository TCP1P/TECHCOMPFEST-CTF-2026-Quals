from Crypto.Cipher import AES
from Crypto.Util.Padding import *
import random
import subprocess
import json
import os
import shutil
import tempfile
import re
from typing import Tuple, Optional, Any

class SecureExifToolWrapper:
    def __init__(self, executable: str = "exiftool"):
        self.executable = shutil.which(executable)
        if not self.executable:
            raise FileNotFoundError("ExifTool binary not found.")

        # Layer 1: Whitelist of allowed functional flags
        self.ALLOWED_FLAGS = {
            "-all", "-G", "-n", "-s", "-h", "-X", "-t", "-m", "-q", 
            "-validate", "-warning", "-b", "-struct", "-j", "-csv", 
            "-p", "-f", "-L", "-charset", "-a", "-u", "-U", "-H", 
            "-groupNames", "-fast", "-ext", "-r", "-P", 
            "-overwrite_original", "-overwrite_original_in_place", 
            "-d", "-coordFormat", "-sep", "-config", "-execute", 
            "-common_args", "-v", "-ee", "-api", "-lang", "--"
        }
        
        # Regex for Layer 2: Only alphanumeric and underscores
        self.safe_pattern = re.compile(r"^[a-zA-Z0-9_]*$")

    def _is_safe(self, arg: str) -> bool:
        """Validates that an argument is either a whitelisted flag or a safe tag/value."""
        # Check if it's a whitelisted functional flag
        if arg in self.ALLOWED_FLAGS:
            return True
        
        # validate value
        if "=" in arg:
            tag, value = arg.split("=", 1)
            if tag in self.ALLOWED_FLAGS:
                return bool(self.safe_pattern.match(value))
            else:
                return False
            
        # Otherwise, check if the value is strictly alphanumeric
        return bool(self.safe_pattern.match(arg))

    def run(self, *args: str, filepath: str) -> Tuple[Optional[Any], str]:
        """Runs exiftool on a copy with strict argument validation."""
        if not os.path.exists(filepath):
            return None, f"Error: File '{filepath}' not found."

        # Validate all arguments before building the command
        sanitized_args = []
        for arg in args:
            if self._is_safe(arg):
                sanitized_args.append(arg)
            else:
                return None, f"Security Block: Illegal argument detected -> '{arg}'"

        with tempfile.TemporaryDirectory() as tmpdir:
            original_filename = os.path.basename(filepath)
            temp_filepath = os.path.join(tmpdir, original_filename)
            shutil.copy2(filepath, temp_filepath)

            command = [self.executable] + sanitized_args + [temp_filepath]
            read = [self.executable, "-j"] + [temp_filepath]

            try:
                # Using list-based Popen is safer than shell=True, 
                # but validation adds the necessary second layer.
                process = subprocess.Popen(
                    command,
                    stdout=subprocess.PIPE,
                    stderr=subprocess.PIPE,
                    text=True
                )
                
                stdout, stderr = process.communicate()
                
                json_data = None
                
                proc_read = subprocess.Popen(
                    read,
                    stdout=subprocess.PIPE,
                    stderr=subprocess.PIPE,
                    text=True
                )
                
                stdout, _ = proc_read.communicate()
                
                if stdout.strip():
                    try:
                        json_data = json.loads(stdout)
                        if isinstance(json_data, list) and len(json_data) > 0:
                            json_data[0]["SourceFile"] = filepath
                    except json.JSONDecodeError:
                        stderr += "\nInternal Error: Failed to parse JSON."
                return json_data, stderr

            except Exception as e:
                return None, f"Execution Error: {str(e)}"

tool = SecureExifToolWrapper()

print("Exifdition 33....")
seed = b"GOLDSHIP"

def gen():
    t = list(seed)
    random.shuffle(t)
    return bytearray(t)

k1 = gen()*2
iv1 = gen()*2
k2 = gen()*2
iv2 = gen()*2

def output(msg: str):
    cipher = AES.new(k1, AES.MODE_CBC, iv=iv1)
    ct = cipher.encrypt(pad(msg.encode(), 16))
    cipher = AES.new(k2, AES.MODE_CBC, iv=iv2)
    ct = cipher.encrypt(ct)
    print(ct.hex())

for _ in range(3):
    choice = int(input("?> "))
    
    if choice == 1:
        data, err = tool.run("-j", filepath="flag.png")
        if data:
            output(str(data)+err)
        else:
            print("Error")
    elif choice == 2:
        cmd = input("Command: ")
        data, err = tool.run(cmd, filepath="images.jpg")
        if data:
            output(str(data)+err)
        else:
            print("Error")
    else:
        print("Bye")
        exit(1)
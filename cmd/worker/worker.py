import sys
import base64
import cv2
import numpy as np
from tensorflow.keras.models import load_model
from tensorflow.keras.losses import binary_crossentropy


def split(x, num_or_size_splits, axis=0):
    return np.array_split(x, num_or_size_splits, axis=axis)


def reduce_sum(x, axis=None, keepdims=False):
    return np.sum(x, axis=axis, keepdims=keepdims)


def dice_coef(y_true, y_pred):
    y_pred_split = split(y_pred, 3, axis=-1)
    y_true_split = split(y_true, 3, axis=-1)
    dice_summ = 0
    for a_y_pred, b_y_true in zip(y_pred_split, y_true_split):
        dice_calculate = (2 * reduce_sum(a_y_pred * b_y_true) + 1) / (
            reduce_sum(a_y_pred + b_y_true) + 1
        )
        dice_summ += dice_calculate
    return dice_summ / 3


def dice_loss(y_true, y_pred):
    return 1 - dice_coef(y_true, y_pred)


def loss_func(y_true, y_pred):
    return 0.3 * dice_loss(y_true, y_pred) + binary_crossentropy(y_true, y_pred)


def base64_to_image(base64_string, target_size=(256, 256)):

    image_data = base64.b64decode(base64_string)
    nparr = np.frombuffer(image_data, np.uint8)
    img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)

    if img is None:
        raise ValueError("Failed to decode Base64 image")

    img = cv2.resize(img, target_size)

    img_batch = np.expand_dims(img, axis=0)

    return img_batch


def colored_mask_to_base64(predicted_mask):
 
    colored_mask = (predicted_mask[0] * 255).astype(np.uint8)

    _, buffer = cv2.imencode('.png', colored_mask)
    base64_mask = base64.b64encode(buffer).decode('utf-8')

    return base64_mask



if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Error: input data required", file=sys.stderr)
        sys.exit(1)

    try:
        img_batch = base64_to_image(sys.argv[1])

        Segmentator = load_model(
            "/segmentator_weights_converted.h5",
            custom_objects={
                "dice_coef": dice_coef,
                "dice_loss": dice_loss,
                "loss_func": loss_func,
            },
        )

        predicted_mask = Segmentator.predict(img_batch, verbose=0)
        result = colored_mask_to_base64(predicted_mask)

        sys.stdout.write(result)
        #sys.stdout.flush()
        
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)
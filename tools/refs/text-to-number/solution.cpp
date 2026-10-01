#include <climits>
#include <string>
using namespace std;

int textToNumber(string& text) {
    size_t i = 0;
    while (i < text.size() && text[i] == ' ') i++;
    int sign = 1;
    if (i < text.size() && (text[i] == '+' || text[i] == '-')) sign = text[i++] == '-' ? -1 : 1;
    long long value = 0;
    for (; i < text.size() && isdigit((unsigned char)text[i]); i++) {
        value = value * 10 + (text[i] - '0');
        if (value > INT_MAX) return sign == 1 ? INT_MAX : INT_MIN;
    }
    return sign * value;
}

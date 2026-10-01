#include <string>
using namespace std;

string compress(string& word) {
    string out;
    for (size_t i = 0; i < word.size();) {
        size_t j = i;
        while (j < word.size() && word[j] == word[i]) j++;
        out += word[i];
        if (j - i > 1) out += to_string(j - i);
        i = j;
    }
    return out;
}

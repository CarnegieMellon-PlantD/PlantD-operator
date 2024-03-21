# Gofakeit Functions

## Payment

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>achaccount</code></td>
<td>A bank account number used for Automated Clearing House transactions and electronic transfers</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>achrouting</code></td>
<td>Unique nine-digit code used in the U.S. for identifying the bank and processing electronic transactions</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>bitcoinaddress</code></td>
<td>Cryptographic identifier used to receive, store, and send Bitcoin cryptocurrency in a peer-to-peer network</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>bitcoinprivatekey</code></td>
<td>Secret, secure code that allows the owner to access and control their Bitcoin holdings</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>creditcard</code></td>
<td>Plastic card allowing users to make purchases on credit, with payment due at a later date</td>
<td><code>map[string]any</code></td>
<td></td>
</tr>
<tr>
<td><code>creditcardcvv</code></td>
<td>Three or four-digit security code on a credit card used for online and remote transactions</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>creditcardexp</code></td>
<td>Date when a credit card becomes invalid and cannot be used for transactions</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>creditcardnumber</code></td>
<td>Unique numerical identifier on a credit card used for making electronic payments and transactions</td>
<td><code>string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>types</code></td>
<td>A select number of types you want to use when generating a credit card number</td>
<td><code>[]string</code></td>
<td>False</td>
<td><code>all</code></td>
<td>
<li><code>visa</code></li>
<li><code>mastercard</code></li>
<li><code>american-express</code></li>
<li><code>diners-club</code></li>
<li><code>discover</code></li>
<li><code>jcb</code></li>
<li><code>unionpay</code></li>
<li><code>maestro</code></li>
<li><code>elo</code></li>
<li><code>hiper</code></li>
<li><code>hipercard</code></li>
</td>
</tr>
<tr>
<td><code>bins</code></td>
<td>Optional list of prepended bin numbers to pick from</td>
<td><code>[]string</code></td>
<td>True</td>
<td></td>
<td></td>
</tr>
<tr>
<td><code>gaps</code></td>
<td>Whether or not to have gaps in number</td>
<td><code>bool</code></td>
<td>False</td>
<td><code>false</code></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>creditcardtype</code></td>
<td>Classification of credit cards based on the issuing company</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>currency</code></td>
<td>Medium of exchange, often in the form of paper money or coins, used for trade and transactions</td>
<td><code>map[string]string</code></td>
<td></td>
</tr>
<tr>
<td><code>currencylong</code></td>
<td>Complete name of a specific currency used for official identification in financial transactions</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>currencyshort</code></td>
<td>Short 3-letter word used to represent a specific currency</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>price</code></td>
<td>The amount of money or value assigned to a product, service, or asset in a transaction</td>
<td><code>float64</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>min</code></td>
<td>Minimum price value</td>
<td><code>float</code></td>
<td>False</td>
<td><code>0</code></td>
<td></td>
</tr>
<tr>
<td><code>max</code></td>
<td>Maximum price value</td>
<td><code>float</code></td>
<td>False</td>
<td><code>1000</code></td>
<td></td>
</tr>
</table></td>
</tr>
</table>

# Gofakeit Functions

## Address

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>address</code></td>
<td>Residential location including street, city, state, country and postal code</td>
<td><code>map[string]any</code></td>
<td></td>
</tr>
<tr>
<td><code>city</code></td>
<td>Part of a country with significant population, often a central hub for culture and commerce</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>country</code></td>
<td>Nation with its own government and defined territory</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>countryabr</code></td>
<td>Shortened 2-letter form of a country's name</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>latitude</code></td>
<td>Geographic coordinate specifying north-south position on Earth's surface</td>
<td><code>float</code></td>
<td></td>
</tr>
<tr>
<td><code>latituderange</code></td>
<td>Latitude number between the given range (default min=0, max=90)</td>
<td><code>float</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>min</code></td>
<td>Minimum range</td>
<td><code>float</code></td>
<td>False</td>
<td><code>0</code></td>
<td></td>
</tr>
<tr>
<td><code>max</code></td>
<td>Maximum range</td>
<td><code>float</code></td>
<td>False</td>
<td><code>90</code></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>longitude</code></td>
<td>Geographic coordinate indicating east-west position on Earth's surface</td>
<td><code>float</code></td>
<td></td>
</tr>
<tr>
<td><code>longituderange</code></td>
<td>Longitude number between the given range (default min=0, max=180)</td>
<td><code>float</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>min</code></td>
<td>Minimum range</td>
<td><code>float</code></td>
<td>False</td>
<td><code>0</code></td>
<td></td>
</tr>
<tr>
<td><code>max</code></td>
<td>Maximum range</td>
<td><code>float</code></td>
<td>False</td>
<td><code>180</code></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>state</code></td>
<td>Governmental division within a country, often having its own laws and government</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>stateabr</code></td>
<td>Shortened 2-letter form of a country's state</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>street</code></td>
<td>Public road in a city or town, typically with houses and buildings on each side</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>streetname</code></td>
<td>Name given to a specific road or street</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>streetnumber</code></td>
<td>Numerical identifier assigned to a street</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>streetprefix</code></td>
<td>Directional or descriptive term preceding a street name, like 'East' or 'Main'</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>streetsuffix</code></td>
<td>Designation at the end of a street name indicating type, like 'Avenue' or 'Street'</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>zip</code></td>
<td>Numerical code for postal address sorting, specific to a geographic area</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Word

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>adjective</code></td>
<td>Word describing or modifying a noun</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>adjectivedemonstrative</code></td>
<td>Adjective used to point out specific things</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>adjectivedescriptive</code></td>
<td>Adjective that provides detailed characteristics about a noun</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>adjectiveindefinite</code></td>
<td>Adjective describing a non-specific noun</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>adjectiveinterrogative</code></td>
<td>Adjective used to ask questions</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>adjectivepossessive</code></td>
<td>Adjective indicating ownership or possession</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>adjectiveproper</code></td>
<td>Adjective derived from a proper noun, often used to describe nationality or origin</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>adjectivequantitative</code></td>
<td>Adjective that indicates the quantity or amount of something</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>adverb</code></td>
<td>Word that modifies verbs, adjectives, or other adverbs</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>adverbdegree</code></td>
<td>Adverb that indicates the degree or intensity of an action or adjective</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>adverbfrequencydefinite</code></td>
<td>Adverb that specifies how often an action occurs with a clear frequency</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>adverbfrequencyindefinite</code></td>
<td>Adverb that specifies how often an action occurs without specifying a particular frequency</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>adverbmanner</code></td>
<td>Adverb that describes how an action is performed</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>adverbplace</code></td>
<td>Adverb that indicates the location or direction of an action</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>adverbtimedefinite</code></td>
<td>Adverb that specifies the exact time an action occurs</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>adverbtimeindefinite</code></td>
<td>Adverb that gives a general or unspecified time frame</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>comment</code></td>
<td>Statement or remark expressing an opinion, observation, or reaction</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>connective</code></td>
<td>Word used to connect words or sentences</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>connectivecasual</code></td>
<td>Connective word used to indicate a cause-and-effect relationship between events or actions</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>connectivecomparative</code></td>
<td>Connective word used to indicate a comparison between two or more things</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>connectivecomplaint</code></td>
<td>Connective word used to express dissatisfaction or complaints about a situation</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>connectiveexamplify</code></td>
<td>Connective word used to provide examples or illustrations of a concept or idea</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>connectivelisting</code></td>
<td>Connective word used to list or enumerate items or examples</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>connectivetime</code></td>
<td>Connective word used to indicate a temporal relationship between events or actions</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>interjection</code></td>
<td>Word expressing emotion</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>loremipsumparagraph</code></td>
<td>Paragraph of the Lorem Ipsum placeholder text used in design and publishing</td>
<td><code>string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>paragraphcount</code></td>
<td>Number of paragraphs</td>
<td><code>int</code></td>
<td>False</td>
<td><code>2</code></td>
<td></td>
</tr>
<tr>
<td><code>sentencecount</code></td>
<td>Number of sentences in a paragraph</td>
<td><code>int</code></td>
<td>False</td>
<td><code>2</code></td>
<td></td>
</tr>
<tr>
<td><code>wordcount</code></td>
<td>Number of words in a sentence</td>
<td><code>int</code></td>
<td>False</td>
<td><code>5</code></td>
<td></td>
</tr>
<tr>
<td><code>paragraphseparator</code></td>
<td>String value to add between paragraphs</td>
<td><code>string</code></td>
<td>False</td>
<td><code>&lt;br /&gt;</code></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>loremipsumsentence</code></td>
<td>Sentence of the Lorem Ipsum placeholder text used in design and publishing</td>
<td><code>string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>wordcount</code></td>
<td>Number of words in a sentence</td>
<td><code>int</code></td>
<td>False</td>
<td><code>5</code></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>loremipsumword</code></td>
<td>Word of the Lorem Ipsum placeholder text used in design and publishing</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>noun</code></td>
<td>Person, place, thing, or idea, named or referred to in a sentence</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>nounabstract</code></td>
<td>Ideas, qualities, or states that cannot be perceived with the five senses</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>nouncollectiveanimal</code></td>
<td>Group of animals, like a 'pack' of wolves or a 'flock' of birds</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>nouncollectivepeople</code></td>
<td>Group of people or things regarded as a unit</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>nouncollectivething</code></td>
<td>Group of objects or items, such as a 'bundle' of sticks or a 'cluster' of grapes</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>nouncommon</code></td>
<td>General name for people, places, or things, not specific or unique</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>nounconcrete</code></td>
<td>Names for physical entities experienced through senses like sight, touch, smell, or taste</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>nouncountable</code></td>
<td>Items that can be counted individually</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>noundeterminer</code></td>
<td>Word that introduces a noun and identifies it as a noun</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>nounproper</code></td>
<td>Specific name for a particular person, place, or organization</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>noununcountable</code></td>
<td>Items that can't be counted individually</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>paragraph</code></td>
<td>Distinct section of writing covering a single theme, composed of multiple sentences</td>
<td><code>string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>paragraphcount</code></td>
<td>Number of paragraphs</td>
<td><code>int</code></td>
<td>False</td>
<td><code>2</code></td>
<td></td>
</tr>
<tr>
<td><code>sentencecount</code></td>
<td>Number of sentences in a paragraph</td>
<td><code>int</code></td>
<td>False</td>
<td><code>2</code></td>
<td></td>
</tr>
<tr>
<td><code>wordcount</code></td>
<td>Number of words in a sentence</td>
<td><code>int</code></td>
<td>False</td>
<td><code>5</code></td>
<td></td>
</tr>
<tr>
<td><code>paragraphseparator</code></td>
<td>String value to add between paragraphs</td>
<td><code>string</code></td>
<td>False</td>
<td><code>&lt;br /&gt;</code></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>phrase</code></td>
<td>A small group of words standing together</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>phraseadverb</code></td>
<td>Phrase that modifies a verb, adjective, or another adverb, providing additional information.</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>phrasenoun</code></td>
<td>Phrase with a noun as its head, functions within sentence like a noun</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>phrasepreposition</code></td>
<td>Phrase starting with a preposition, showing relation between elements in a sentence.</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>phraseverb</code></td>
<td>Phrase that Consists of a verb and its modifiers, expressing an action or state</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>preposition</code></td>
<td>Words used to express the relationship of a noun or pronoun to other words in a sentence</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>prepositioncompound</code></td>
<td>Preposition that can be formed by combining two or more prepositions</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>prepositiondouble</code></td>
<td>Two-word combination preposition, indicating a complex relation</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>prepositionsimple</code></td>
<td>Single-word preposition showing relationships between 2 parts of a sentence</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>pronoun</code></td>
<td>Word used in place of a noun to avoid repetition</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>pronoundemonstrative</code></td>
<td>Pronoun that points out specific people or things</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>pronounindefinite</code></td>
<td>Pronoun that does not refer to a specific person or thing</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>pronouninterrogative</code></td>
<td>Pronoun used to ask questions</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>pronounobject</code></td>
<td>Pronoun used as the object of a verb or preposition</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>pronounpersonal</code></td>
<td>Pronoun referring to a specific persons or things</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>pronounpossessive</code></td>
<td>Pronoun indicating ownership or belonging</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>pronounreflective</code></td>
<td>Pronoun referring back to the subject of the sentence</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>pronounrelative</code></td>
<td>Pronoun that introduces a clause, referring back to a noun or pronoun</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>question</code></td>
<td>Statement formulated to inquire or seek clarification</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>quote</code></td>
<td>Direct repetition of someone else's words</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>sentence</code></td>
<td>Set of words expressing a statement, question, exclamation, or command</td>
<td><code>string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>wordcount</code></td>
<td>Number of words in a sentence</td>
<td><code>int</code></td>
<td>False</td>
<td><code>5</code></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>sentencesimple</code></td>
<td>Group of words that expresses a complete thought</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>verb</code></td>
<td>Word expressing an action, event or state</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>verbaction</code></td>
<td>Verb Indicating a physical or mental action</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>verbhelping</code></td>
<td>Auxiliary verb that helps the main verb complete the sentence</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>verbintransitive</code></td>
<td>Verb that does not require a direct object to complete its meaning</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>verblinking</code></td>
<td>Verb that Connects the subject of a sentence to a subject complement</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>verbtransitive</code></td>
<td>Verb that requires a direct object to complete its meaning</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>word</code></td>
<td>Basic unit of language representing a concept or thing, consisting of letters and having meaning</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Animal

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>animal</code></td>
<td>Living creature with the ability to move, eat, and interact with its environment</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>animaltype</code></td>
<td>Type of animal, such as mammals, birds, reptiles, etc.</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>bird</code></td>
<td>Distinct species of birds</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>cat</code></td>
<td>Various breeds that define different cats</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>dog</code></td>
<td>Various breeds that define different dogs</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>farmanimal</code></td>
<td>Animal name commonly found on a farm</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>petname</code></td>
<td>Affectionate nickname given to a pet</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## App

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>appauthor</code></td>
<td>Person or group creating and developing an application</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>appname</code></td>
<td>Software program designed for a specific purpose or task on a computer or mobile device</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>appversion</code></td>
<td>Particular release of an application in Semantic Versioning format</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Beer

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>beeralcohol</code></td>
<td>Measures the alcohol content in beer</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>beerblg</code></td>
<td>Scale indicating the concentration of extract in worts</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>beerhop</code></td>
<td>The flower used in brewing to add flavor, aroma, and bitterness to beer</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>beeribu</code></td>
<td>Scale measuring bitterness of beer from hops</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>beermalt</code></td>
<td>Processed barley or other grains, provides sugars for fermentation and flavor to beer</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>beername</code></td>
<td>Specific brand or variety of beer</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>beerstyle</code></td>
<td>Distinct characteristics and flavors of beer</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>beeryeast</code></td>
<td>Microorganism used in brewing to ferment sugars, producing alcohol and carbonation in beer</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Company

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>blurb</code></td>
<td>Brief description or summary of a company's purpose, products, or services</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>bs</code></td>
<td>Random bs company word</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>buzzword</code></td>
<td>Trendy or overused term often used in business to sound impressive</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>company</code></td>
<td>Designated official name of a business or organization</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>companysuffix</code></td>
<td>Suffix at the end of a company name, indicating business structure, like 'Inc.' or 'LLC'</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>job</code></td>
<td>Position or role in employment, involving specific tasks and responsibilities</td>
<td><code>map[string]string</code></td>
<td></td>
</tr>
<tr>
<td><code>jobdescriptor</code></td>
<td>Word used to describe the duties, requirements, and nature of a job</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>joblevel</code></td>
<td>Random job level</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>jobtitle</code></td>
<td>Specific title for a position or role within a company or organization</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>slogan</code></td>
<td>Catchphrase or motto used by a company to represent its brand or values</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Book

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>book</code></td>
<td>Written or printed work consisting of pages bound together, covering various subjects or stories</td>
<td><code>map[string]string</code></td>
<td></td>
</tr>
<tr>
<td><code>bookauthor</code></td>
<td>The individual who wrote or created the content of a book</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>bookgenre</code></td>
<td>Category or type of book defined by its content, style, or form</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>booktitle</code></td>
<td>The specific name given to a book</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Misc

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>bool</code></td>
<td>Data type that represents one of two possible values, typically true or false</td>
<td><code>bool</code></td>
<td></td>
</tr>
<tr>
<td><code>flipacoin</code></td>
<td>Decision-making method involving the tossing of a coin to determine outcomes</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>uuid</code></td>
<td>128-bit identifier used to uniquely identify objects or entities in computer systems</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>weighted</code></td>
<td>Randomly select a given option based upon an equal amount of weights</td>
<td><code>any</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>options</code></td>
<td>Array of any values</td>
<td><code>[]string</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
<tr>
<td><code>weights</code></td>
<td>Array of weights</td>
<td><code>[]float</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
</table>

# Gofakeit Functions

## Food

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>breakfast</code></td>
<td>First meal of the day, typically eaten in the morning</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>dessert</code></td>
<td>Sweet treat often enjoyed after a meal</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>dinner</code></td>
<td>Evening meal, typically the day's main and most substantial meal</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>drink</code></td>
<td>Liquid consumed for hydration, pleasure, or nutritional benefits</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>fruit</code></td>
<td>Edible plant part, typically sweet, enjoyed as a natural snack or dessert</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>lunch</code></td>
<td>Midday meal, often lighter than dinner, eaten around noon</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>snack</code></td>
<td>Random snack</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>vegetable</code></td>
<td>Edible plant or part of a plant, often used in savory cooking or salads</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Car

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>car</code></td>
<td>Wheeled motor vehicle used for transportation</td>
<td><code>map[string]any</code></td>
<td></td>
</tr>
<tr>
<td><code>carfueltype</code></td>
<td>Type of energy source a car uses</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>carmaker</code></td>
<td>Company or brand that manufactures and designs cars</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>carmodel</code></td>
<td>Specific design or version of a car produced by a manufacturer</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>cartransmissiontype</code></td>
<td>Mechanism a car uses to transmit power from the engine to the wheels</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>cartype</code></td>
<td>Classification of cars based on size, use, or body style</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Celebrity

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>celebrityactor</code></td>
<td>Famous person known for acting in films, television, or theater</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>celebritybusiness</code></td>
<td>High-profile individual known for significant achievements in business or entrepreneurship</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>celebritysport</code></td>
<td>Famous athlete known for achievements in a particular sport</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Internet

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>chromeuseragent</code></td>
<td>The specific identification string sent by the Google Chrome web browser when making requests on the internet</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>domainname</code></td>
<td>Human-readable web address used to identify websites on the internet</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>domainsuffix</code></td>
<td>The part of a domain name that comes after the last dot, indicating its type or purpose</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>firefoxuseragent</code></td>
<td>The specific identification string sent by the Firefox web browser when making requests on the internet</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>httpmethod</code></td>
<td>Verb used in HTTP requests to specify the desired action to be performed on a resource</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>httpstatuscode</code></td>
<td>Random http status code</td>
<td><code>int</code></td>
<td></td>
</tr>
<tr>
<td><code>httpstatuscodesimple</code></td>
<td>Three-digit number returned by a web server to indicate the outcome of an HTTP request</td>
<td><code>int</code></td>
<td></td>
</tr>
<tr>
<td><code>httpversion</code></td>
<td>Number indicating the version of the HTTP protocol used for communication between a client and a server</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>ipv4address</code></td>
<td>Numerical label assigned to devices on a network for identification and communication</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>ipv6address</code></td>
<td>Numerical label assigned to devices on a network, providing a larger address space than IPv4 for internet communication</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>loglevel</code></td>
<td>Classification used in logging to indicate the severity or priority of a log entry</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>macaddress</code></td>
<td>Unique identifier assigned to network interfaces, often used in Ethernet networks</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>operauseragent</code></td>
<td>The specific identification string sent by the Opera web browser when making requests on the internet</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>safariuseragent</code></td>
<td>The specific identification string sent by the Safari web browser when making requests on the internet</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>url</code></td>
<td>Web address that specifies the location of a resource on the internet</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>useragent</code></td>
<td>String sent by a web browser to identify itself when requesting web content</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Color

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>color</code></td>
<td>Hue seen by the eye, returns the name of the color like red or blue</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>hexcolor</code></td>
<td>Six-digit code representing a color in the color model</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>nicecolors</code></td>
<td>Attractive and appealing combinations of colors, returns an list of color hex codes</td>
<td><code>[]string</code></td>
<td></td>
</tr>
<tr>
<td><code>rgbcolor</code></td>
<td>Color defined by red, green, and blue light values</td>
<td><code>[]int</code></td>
<td></td>
</tr>
<tr>
<td><code>safecolor</code></td>
<td>Colors displayed consistently on different web browsers and devices</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## File

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>csv</code></td>
<td>Individual lines or data entries within a CSV (Comma-Separated Values) format</td>
<td><code>[]byte</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>delimiter</code></td>
<td>Separator in between row values</td>
<td><code>string</code></td>
<td>False</td>
<td><code>,</code></td>
<td></td>
</tr>
<tr>
<td><code>rowcount</code></td>
<td>Number of rows</td>
<td><code>int</code></td>
<td>False</td>
<td><code>100</code></td>
<td></td>
</tr>
<tr>
<td><code>fields</code></td>
<td>Fields containing key name and function</td>
<td><code>[]Field</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>fileextension</code></td>
<td>Suffix appended to a filename indicating its format or type</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>filemimetype</code></td>
<td>Defines file format and nature for browsers and email clients using standardized identifiers</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>json</code></td>
<td>Format for structured data interchange used in programming, returns an object or an array of objects</td>
<td><code>[]byte</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>type</code></td>
<td>Type of JSON, object or array</td>
<td><code>string</code></td>
<td>False</td>
<td><code>object</code></td>
<td>
<li><code>object</code></li>
<li><code>array</code></li>
</td>
</tr>
<tr>
<td><code>rowcount</code></td>
<td>Number of rows in JSON array</td>
<td><code>int</code></td>
<td>False</td>
<td><code>100</code></td>
<td></td>
</tr>
<tr>
<td><code>indent</code></td>
<td>Whether or not to add indents and newlines</td>
<td><code>bool</code></td>
<td>False</td>
<td><code>false</code></td>
<td></td>
</tr>
<tr>
<td><code>fields</code></td>
<td>Fields containing key name and function to run in json format</td>
<td><code>[]Field</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>xml</code></td>
<td>Generates an single or an array of elements in xml format</td>
<td><code>[]byte</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>type</code></td>
<td>Type of XML, single or array</td>
<td><code>string</code></td>
<td>False</td>
<td><code>single</code></td>
<td>
<li><code>single</code></li>
<li><code>array</code></li>
</td>
</tr>
<tr>
<td><code>rootelement</code></td>
<td>Root element wrapper name</td>
<td><code>string</code></td>
<td>False</td>
<td><code>xml</code></td>
<td></td>
</tr>
<tr>
<td><code>recordelement</code></td>
<td>Record element for each record row</td>
<td><code>string</code></td>
<td>False</td>
<td><code>record</code></td>
<td></td>
</tr>
<tr>
<td><code>rowcount</code></td>
<td>Number of rows in JSON array</td>
<td><code>int</code></td>
<td>False</td>
<td><code>100</code></td>
<td></td>
</tr>
<tr>
<td><code>indent</code></td>
<td>Whether or not to add indents and newlines</td>
<td><code>bool</code></td>
<td>False</td>
<td><code>false</code></td>
<td></td>
</tr>
<tr>
<td><code>fields</code></td>
<td>Fields containing key name and function to run in json format</td>
<td><code>[]Field</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
</table>

# Gofakeit Functions

## Finance

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>cusip</code></td>
<td>Unique identifier for securities, especially bonds, in the United States and Canada</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>isin</code></td>
<td>International standard code for uniquely identifying securities worldwide</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Time

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>date</code></td>
<td>Representation of a specific day, month, and year, often used for chronological reference</td>
<td><code>string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>format</code></td>
<td>Date time string format output. You may also use golang time format or java time format</td>
<td><code>string</code></td>
<td>False</td>
<td><code>RFC3339</code></td>
<td>
<li><code>ANSIC</code></li>
<li><code>UnixDate</code></li>
<li><code>RubyDate</code></li>
<li><code>RFC822</code></li>
<li><code>RFC822Z</code></li>
<li><code>RFC850</code></li>
<li><code>RFC1123</code></li>
<li><code>RFC1123Z</code></li>
<li><code>RFC3339</code></li>
<li><code>RFC3339Nano</code></li>
</td>
</tr>
</table></td>
</tr>
<tr>
<td><code>daterange</code></td>
<td>Random date between two ranges</td>
<td><code>string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>startdate</code></td>
<td>Start date time string</td>
<td><code>string</code></td>
<td>False</td>
<td><code>1970-01-01</code></td>
<td></td>
</tr>
<tr>
<td><code>enddate</code></td>
<td>End date time string</td>
<td><code>string</code></td>
<td>False</td>
<td><code>2024-03-21</code></td>
<td></td>
</tr>
<tr>
<td><code>format</code></td>
<td>Date time string format</td>
<td><code>string</code></td>
<td>False</td>
<td><code>yyyy-MM-dd</code></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>day</code></td>
<td>24-hour period equivalent to one rotation of Earth on its axis</td>
<td><code>int</code></td>
<td></td>
</tr>
<tr>
<td><code>futuretime</code></td>
<td>Date that has occurred after the current moment in time</td>
<td><code>time</code></td>
<td></td>
</tr>
<tr>
<td><code>hour</code></td>
<td>Unit of time equal to 60 minutes</td>
<td><code>int</code></td>
<td></td>
</tr>
<tr>
<td><code>minute</code></td>
<td>Unit of time equal to 60 seconds</td>
<td><code>int</code></td>
<td></td>
</tr>
<tr>
<td><code>month</code></td>
<td>Division of the year, typically 30 or 31 days long</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>monthstring</code></td>
<td>String Representation of a month name</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>nanosecond</code></td>
<td>Unit of time equal to One billionth (10^-9) of a second</td>
<td><code>int</code></td>
<td></td>
</tr>
<tr>
<td><code>pasttime</code></td>
<td>Date that has occurred before the current moment in time</td>
<td><code>time</code></td>
<td></td>
</tr>
<tr>
<td><code>second</code></td>
<td>Unit of time equal to 1/60th of a minute</td>
<td><code>int</code></td>
<td></td>
</tr>
<tr>
<td><code>timezone</code></td>
<td>Region where the same standard time is used, based on longitudinal divisions of the Earth</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>timezoneabv</code></td>
<td>Abbreviated 3-letter word of a timezone</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>timezonefull</code></td>
<td>Full name of a timezone</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>timezoneoffset</code></td>
<td>The difference in hours from Coordinated Universal Time (UTC) for a specific region</td>
<td><code>float32</code></td>
<td></td>
</tr>
<tr>
<td><code>timezoneregion</code></td>
<td>Geographic area sharing the same standard time</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>weekday</code></td>
<td>Day of the week excluding the weekend</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>year</code></td>
<td>Period of 365 days, the time Earth takes to orbit the Sun</td>
<td><code>int</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Game

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>dice</code></td>
<td>Small, cube-shaped objects used in games of chance for random outcomes</td>
<td><code>[]uint</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>numdice</code></td>
<td>Number of dice to roll</td>
<td><code>uint</code></td>
<td>False</td>
<td><code>1</code></td>
<td></td>
</tr>
<tr>
<td><code>sides</code></td>
<td>Number of sides on each dice</td>
<td><code>[]uint</code></td>
<td>False</td>
<td><code>[6]</code></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>gamertag</code></td>
<td>User-selected online username or alias used for identification in games</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## String

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>digit</code></td>
<td>Numerical symbol used to represent numbers</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>digitn</code></td>
<td>string of length N consisting of ASCII digits</td>
<td><code>string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>count</code></td>
<td>Number of digits to generate</td>
<td><code>uint</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>letter</code></td>
<td>Character or symbol from the American Standard Code for Information Interchange (ASCII) character set</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>lettern</code></td>
<td>ASCII string with length N</td>
<td><code>string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>count</code></td>
<td>Number of digits to generate</td>
<td><code>uint</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>lexify</code></td>
<td>Replace ? with random generated letters</td>
<td><code>string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>str</code></td>
<td>String value to replace ?'s</td>
<td><code>string</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>numerify</code></td>
<td>Replace # with random numerical values</td>
<td><code>string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>str</code></td>
<td>String value to replace #'s</td>
<td><code>string</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>randomstring</code></td>
<td>Return a random string from a string array</td>
<td><code>[]string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>strs</code></td>
<td>Delimited separated strings</td>
<td><code>[]string</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>shufflestrings</code></td>
<td>Shuffle an array of strings</td>
<td><code>[]string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>strs</code></td>
<td>Delimited separated strings</td>
<td><code>[]string</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>vowel</code></td>
<td>Speech sound produced with an open vocal tract</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Person

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>email</code></td>
<td>Electronic mail used for sending digital messages and communication over the internet</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>firstname</code></td>
<td>The name given to a person at birth</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>gender</code></td>
<td>Classification based on social and cultural norms that identifies an individual</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>hobby</code></td>
<td>An activity pursued for leisure and pleasure</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>lastname</code></td>
<td>The family name or surname of an individual</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>middlename</code></td>
<td>Name between a person's first name and last name</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>name</code></td>
<td>The given and family name of an individual</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>nameprefix</code></td>
<td>A title or honorific added before a person's name</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>namesuffix</code></td>
<td>A title or designation added after a person's name</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>person</code></td>
<td>Personal data, like name and contact details, used for identification and communication</td>
<td><code>map[string]any</code></td>
<td></td>
</tr>
<tr>
<td><code>phone</code></td>
<td>Numerical sequence used to contact individuals via telephone or mobile devices</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>phoneformatted</code></td>
<td>Formatted phone number of a person</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>ssn</code></td>
<td>Unique nine-digit identifier used for government and financial purposes in the United States</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>teams</code></td>
<td>Randomly split people into teams</td>
<td><code>map[string][]string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>people</code></td>
<td>Array of people</td>
<td><code>[]string</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
<tr>
<td><code>teams</code></td>
<td>Array of teams</td>
<td><code>[]string</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
</table>

# Gofakeit Functions

## Template

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>email_text</code></td>
<td>Written content of an email message, including the sender's message to the recipient</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>markdown</code></td>
<td>Lightweight markup language used for formatting plain text</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>template</code></td>
<td>Generates document from template</td>
<td><code>string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>template</code></td>
<td>Golang template to generate the document from</td>
<td><code>string</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
<tr>
<td><code>data</code></td>
<td>Custom data to pass to the template</td>
<td><code>string</code></td>
<td>True</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
</table>

# Gofakeit Functions

## Emoji

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>emoji</code></td>
<td>Digital symbol expressing feelings or ideas in text messages and online chats</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>emojialias</code></td>
<td>Alternative name or keyword used to represent a specific emoji in text or code</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>emojicategory</code></td>
<td>Group or classification of emojis based on their common theme or use, like 'smileys' or 'animals'</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>emojidescription</code></td>
<td>Brief explanation of the meaning or emotion conveyed by an emoji</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>emojitag</code></td>
<td>Label or keyword associated with an emoji to categorize or search for it easily</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Error

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>error</code></td>
<td>Message displayed by a computer or software when a problem or mistake is encountered</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>errordatabase</code></td>
<td>A problem or issue encountered while accessing or managing a database</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>errorgrpc</code></td>
<td>Communication failure in the high-performance, open-source universal RPC framework</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>errorhttp</code></td>
<td>A problem with a web http request</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>errorhttpclient</code></td>
<td>Failure or issue occurring within a client software that sends requests to web servers</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>errorhttpserver</code></td>
<td>Failure or issue occurring within a server software that recieves requests from clients</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>errorobject</code></td>
<td>Various categories conveying details about encountered errors</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>errorruntime</code></td>
<td>Malfunction occuring during program execution, often causing abrupt termination or unexpected behavior</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>errorvalidation</code></td>
<td>Occurs when input data fails to meet required criteria or format specifications</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Generate

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>fixed_width</code></td>
<td>Fixed width rows of output data based on input fields</td>
<td><code>[]byte</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>rowcount</code></td>
<td>Number of rows</td>
<td><code>int</code></td>
<td>False</td>
<td><code>10</code></td>
<td></td>
</tr>
<tr>
<td><code>fields</code></td>
<td>Fields name, function and params</td>
<td><code>[]Field</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>generate</code></td>
<td>Random string generated from string value based upon available data sets</td>
<td><code>string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>str</code></td>
<td>String value to generate from</td>
<td><code>string</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>map</code></td>
<td>Data structure that stores key-value pairs</td>
<td><code>map[string]any</code></td>
<td></td>
</tr>
<tr>
<td><code>regex</code></td>
<td>Pattern-matching tool used in text processing to search and manipulate strings</td>
<td><code>string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>str</code></td>
<td>Regex RE2 syntax string</td>
<td><code>string</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
</table>

# Gofakeit Functions

## Number

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>float32</code></td>
<td>Data type representing floating-point numbers with 32 bits of precision in computing</td>
<td><code>float32</code></td>
<td></td>
</tr>
<tr>
<td><code>float32range</code></td>
<td>Float32 value between given range</td>
<td><code>float32</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>min</code></td>
<td>Minimum float32 value</td>
<td><code>float</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
<tr>
<td><code>max</code></td>
<td>Maximum float32 value</td>
<td><code>float</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>float64</code></td>
<td>Data type representing floating-point numbers with 64 bits of precision in computing</td>
<td><code>float64</code></td>
<td></td>
</tr>
<tr>
<td><code>float64range</code></td>
<td>Float64 value between given range</td>
<td><code>float64</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>min</code></td>
<td>Minimum float64 value</td>
<td><code>float</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
<tr>
<td><code>max</code></td>
<td>Maximum float64 value</td>
<td><code>float</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>hexuint</code></td>
<td>Hexadecimal representation of an unsigned integer</td>
<td><code>string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>bitSize</code></td>
<td>Bit size of the unsigned integer</td>
<td><code>int</code></td>
<td>False</td>
<td><code>8</code></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>int</code></td>
<td>Signed integer</td>
<td><code>int</code></td>
<td></td>
</tr>
<tr>
<td><code>int16</code></td>
<td>Signed 16-bit integer, capable of representing values from 32,768 to 32,767</td>
<td><code>int16</code></td>
<td></td>
</tr>
<tr>
<td><code>int32</code></td>
<td>Signed 32-bit integer, capable of representing values from -2,147,483,648 to 2,147,483,647</td>
<td><code>int32</code></td>
<td></td>
</tr>
<tr>
<td><code>int64</code></td>
<td>Signed 64-bit integer, capable of representing values from -9,223,372,036,854,775,808 to -9,223,372,036,854,775,807</td>
<td><code>int64</code></td>
<td></td>
</tr>
<tr>
<td><code>int8</code></td>
<td>Signed 8-bit integer, capable of representing values from -128 to 127</td>
<td><code>int8</code></td>
<td></td>
</tr>
<tr>
<td><code>intn</code></td>
<td>Integer value between 0 and n</td>
<td><code>int</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>n</code></td>
<td>Maximum int value</td>
<td><code>int</code></td>
<td>False</td>
<td><code>2147483647</code></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>intrange</code></td>
<td>Integer value between given range</td>
<td><code>int</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>min</code></td>
<td>Minimum int value</td>
<td><code>int</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
<tr>
<td><code>max</code></td>
<td>Maximum int value</td>
<td><code>int</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>number</code></td>
<td>Mathematical concept used for counting, measuring, and expressing quantities or values</td>
<td><code>int</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>min</code></td>
<td>Minimum integer value</td>
<td><code>int</code></td>
<td>False</td>
<td><code>-2147483648</code></td>
<td></td>
</tr>
<tr>
<td><code>max</code></td>
<td>Maximum integer value</td>
<td><code>int</code></td>
<td>False</td>
<td><code>2147483647</code></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>randomint</code></td>
<td>Randomly selected value from a slice of int</td>
<td><code>int</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>ints</code></td>
<td>Delimited separated integers</td>
<td><code>[]int</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>randomuint</code></td>
<td>Randomly selected value from a slice of uint</td>
<td><code>uint</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>uints</code></td>
<td>Delimited separated unsigned integers</td>
<td><code>[]uint</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>shuffleints</code></td>
<td>Shuffles an array of ints</td>
<td><code>[]int</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>ints</code></td>
<td>Delimited separated integers</td>
<td><code>[]int</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>uint</code></td>
<td>Unsigned integer</td>
<td><code>uint</code></td>
<td></td>
</tr>
<tr>
<td><code>uint16</code></td>
<td>Unsigned 16-bit integer, capable of representing values from 0 to 65,535</td>
<td><code>uint16</code></td>
<td></td>
</tr>
<tr>
<td><code>uint32</code></td>
<td>Unsigned 32-bit integer, capable of representing values from 0 to 4,294,967,295</td>
<td><code>uint32</code></td>
<td></td>
</tr>
<tr>
<td><code>uint64</code></td>
<td>Unsigned 64-bit integer, capable of representing values from 0 to 18,446,744,073,709,551,615</td>
<td><code>uint64</code></td>
<td></td>
</tr>
<tr>
<td><code>uint8</code></td>
<td>Unsigned 8-bit integer, capable of representing values from 0 to 255</td>
<td><code>uint8</code></td>
<td></td>
</tr>
<tr>
<td><code>uintn</code></td>
<td>Unsigned integer between 0 and n</td>
<td><code>uint</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>n</code></td>
<td>Maximum uint value</td>
<td><code>uint</code></td>
<td>False</td>
<td><code>4294967295</code></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>uintrange</code></td>
<td>Non-negative integer value between given range</td>
<td><code>uint</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>min</code></td>
<td>Minimum uint value</td>
<td><code>uint</code></td>
<td>False</td>
<td><code>0</code></td>
<td></td>
</tr>
<tr>
<td><code>max</code></td>
<td>Maximum uint value</td>
<td><code>uint</code></td>
<td>False</td>
<td><code>4294967295</code></td>
<td></td>
</tr>
</table></td>
</tr>
</table>

# Gofakeit Functions

## Hacker

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>hackerabbreviation</code></td>
<td>Abbreviations and acronyms commonly used in the hacking and cybersecurity community</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>hackeradjective</code></td>
<td>Adjectives describing terms often associated with hackers and cybersecurity experts</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>hackeringverb</code></td>
<td>Verb describing actions and activities related to hacking, often involving computer systems and security</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>hackernoun</code></td>
<td>Noun representing an element, tool, or concept within the realm of hacking and cybersecurity</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>hackerphrase</code></td>
<td>Informal jargon and slang used in the hacking and cybersecurity community</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>hackerverb</code></td>
<td>Verbs associated with actions and activities in the field of hacking and cybersecurity</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Hipster

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>hipsterparagraph</code></td>
<td>Paragraph showcasing the use of trendy and unconventional vocabulary associated with hipster culture</td>
<td><code>string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>paragraphcount</code></td>
<td>Number of paragraphs</td>
<td><code>int</code></td>
<td>False</td>
<td><code>2</code></td>
<td></td>
</tr>
<tr>
<td><code>sentencecount</code></td>
<td>Number of sentences in a paragraph</td>
<td><code>int</code></td>
<td>False</td>
<td><code>2</code></td>
<td></td>
</tr>
<tr>
<td><code>wordcount</code></td>
<td>Number of words in a sentence</td>
<td><code>int</code></td>
<td>False</td>
<td><code>5</code></td>
<td></td>
</tr>
<tr>
<td><code>paragraphseparator</code></td>
<td>String value to add between paragraphs</td>
<td><code>string</code></td>
<td>False</td>
<td><code>&lt;br /&gt;</code></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>hipstersentence</code></td>
<td>Sentence showcasing the use of trendy and unconventional vocabulary associated with hipster culture</td>
<td><code>string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>wordcount</code></td>
<td>Number of words in a sentence</td>
<td><code>int</code></td>
<td>False</td>
<td><code>5</code></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>hipsterword</code></td>
<td>Trendy and unconventional vocabulary used by hipsters to express unique cultural preferences</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Image

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>imagejpeg</code></td>
<td>Image file format known for its efficient compression and compatibility</td>
<td><code>[]byte</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>width</code></td>
<td>Image width in px</td>
<td><code>int</code></td>
<td>False</td>
<td><code>500</code></td>
<td></td>
</tr>
<tr>
<td><code>height</code></td>
<td>Image height in px</td>
<td><code>int</code></td>
<td>False</td>
<td><code>500</code></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>imagepng</code></td>
<td>Image file format known for its lossless compression and support for transparency</td>
<td><code>[]byte</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>width</code></td>
<td>Image width in px</td>
<td><code>int</code></td>
<td>False</td>
<td><code>500</code></td>
<td></td>
</tr>
<tr>
<td><code>height</code></td>
<td>Image height in px</td>
<td><code>int</code></td>
<td>False</td>
<td><code>500</code></td>
<td></td>
</tr>
</table></td>
</tr>
</table>

# Gofakeit Functions

## Html

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>inputname</code></td>
<td>Attribute used to define the name of an input element in web forms</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>svg</code></td>
<td>Scalable Vector Graphics used to display vector images in web content</td>
<td><code>string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>width</code></td>
<td>Width in px</td>
<td><code>int</code></td>
<td>False</td>
<td><code>500</code></td>
<td></td>
</tr>
<tr>
<td><code>height</code></td>
<td>Height in px</td>
<td><code>int</code></td>
<td>False</td>
<td><code>500</code></td>
<td></td>
</tr>
<tr>
<td><code>type</code></td>
<td>Sub child element type</td>
<td><code>string</code></td>
<td>True</td>
<td></td>
<td>
<li><code>rect</code></li>
<li><code>circle</code></li>
<li><code>ellipse</code></li>
<li><code>line</code></li>
<li><code>polyline</code></li>
<li><code>polygon</code></li>
</td>
</tr>
<tr>
<td><code>colors</code></td>
<td>Hex or RGB array of colors to use</td>
<td><code>[]string</code></td>
<td>True</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
</table>

# Gofakeit Functions

## Language

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>language</code></td>
<td>System of communication using symbols, words, and grammar to convey meaning between individuals</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>languageabbreviation</code></td>
<td>Shortened form of a language's name</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>languagebcp</code></td>
<td>Set of guidelines and standards for identifying and representing languages in computing and internet protocols</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>programminglanguage</code></td>
<td>Formal system of instructions used to create software and perform computational tasks</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Minecraft

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>minecraftanimal</code></td>
<td>Non-hostile creatures in Minecraft, often used for resources and farming</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>minecraftarmorpart</code></td>
<td>Component of an armor set in Minecraft, such as a helmet, chestplate, leggings, or boots</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>minecraftarmortier</code></td>
<td>Classification system for armor sets in Minecraft, indicating their effectiveness and protection level</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>minecraftbiome</code></td>
<td>Distinctive environmental regions in the game, characterized by unique terrain, vegetation, and weather</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>minecraftdye</code></td>
<td>Items used to change the color of various in-game objects</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>minecraftfood</code></td>
<td>Consumable items in Minecraft that provide nourishment to the player character</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>minecraftmobboss</code></td>
<td>Powerful hostile creature in the game, often found in challenging dungeons or structures</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>minecraftmobhostile</code></td>
<td>Aggressive creatures in the game that actively attack players when encountered</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>minecraftmobneutral</code></td>
<td>Creature in the game that only becomes hostile if provoked, typically defending itself when attacked</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>minecraftmobpassive</code></td>
<td>Non-aggressive creatures in the game that do not attack players</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>minecraftore</code></td>
<td>Naturally occurring minerals found in the game Minecraft, used for crafting purposes</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>minecrafttool</code></td>
<td>Items in Minecraft designed for specific tasks, including mining, digging, and building</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>minecraftvillagerjob</code></td>
<td>The profession or occupation assigned to a villager character in the game</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>minecraftvillagerlevel</code></td>
<td>Measure of a villager's experience and proficiency in their assigned job or profession</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>minecraftvillagerstation</code></td>
<td>Designated area or structure in Minecraft where villagers perform their job-related tasks and trading</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>minecraftweapon</code></td>
<td>Tools and items used in Minecraft for combat and defeating hostile mobs</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>minecraftweather</code></td>
<td>Atmospheric conditions in the game that include rain, thunderstorms, and clear skies, affecting gameplay and ambiance</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>minecraftwood</code></td>
<td>Natural resource in Minecraft, used for crafting various items and building structures</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Movie

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>movie</code></td>
<td>A story told through moving pictures and sound</td>
<td><code>map[string]string</code></td>
<td></td>
</tr>
<tr>
<td><code>moviegenre</code></td>
<td>Category that classifies movies based on common themes, styles, and storytelling approaches</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>moviename</code></td>
<td>Title or name of a specific film used for identification and reference</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Auth

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>password</code></td>
<td>Secret word or phrase used to authenticate access to a system or account</td>
<td><code>string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>lower</code></td>
<td>Whether or not to add lower case characters</td>
<td><code>bool</code></td>
<td>False</td>
<td><code>true</code></td>
<td></td>
</tr>
<tr>
<td><code>upper</code></td>
<td>Whether or not to add upper case characters</td>
<td><code>bool</code></td>
<td>False</td>
<td><code>true</code></td>
<td></td>
</tr>
<tr>
<td><code>numeric</code></td>
<td>Whether or not to add numeric characters</td>
<td><code>bool</code></td>
<td>False</td>
<td><code>true</code></td>
<td></td>
</tr>
<tr>
<td><code>special</code></td>
<td>Whether or not to add special characters</td>
<td><code>bool</code></td>
<td>False</td>
<td><code>true</code></td>
<td></td>
</tr>
<tr>
<td><code>space</code></td>
<td>Whether or not to add spaces</td>
<td><code>bool</code></td>
<td>False</td>
<td><code>false</code></td>
<td></td>
</tr>
<tr>
<td><code>length</code></td>
<td>Number of characters in password</td>
<td><code>int</code></td>
<td>False</td>
<td><code>12</code></td>
<td></td>
</tr>
</table></td>
</tr>
<tr>
<td><code>username</code></td>
<td>Unique identifier assigned to a user for accessing an account or system</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Product

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>product</code></td>
<td>An item created for sale or use</td>
<td><code>map[string]any</code></td>
<td></td>
</tr>
<tr>
<td><code>productcategory</code></td>
<td>Classification grouping similar products based on shared characteristics or functions</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>productdescription</code></td>
<td>Explanation detailing the features and characteristics of a product</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>productfeature</code></td>
<td>Specific characteristic of a product that distinguishes it from others products</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>productmaterial</code></td>
<td>The substance from which a product is made, influencing its appearance, durability, and properties</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>productname</code></td>
<td>Distinctive title or label assigned to a product for identification and marketing</td>
<td><code>string</code></td>
<td></td>
</tr>
<tr>
<td><code>productupc</code></td>
<td>Standardized barcode used for product identification and tracking in retail and commerce</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## School

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>school</code></td>
<td>An institution for formal education and learning</td>
<td><code>string</code></td>
<td></td>
</tr>
</table>

# Gofakeit Functions

## Database

<table>
<tr>
<td>Name</td>
<td>Description</td>
<td>Type</td>
<td>Parameters</td>
</tr>
<tr>
<td><code>sql</code></td>
<td>Command in SQL used to add new data records into a database table</td>
<td><code>string</code></td>
<td><table>
<tr>
<th>Name</th>
<th>Description</th>
<th>Type</th>
<th>Optional</th>
<th>Default</th>
<th>Options</th>
</tr>
<tr>
<td><code>table</code></td>
<td>Name of the table to insert into</td>
<td><code>string</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
<tr>
<td><code>count</code></td>
<td>Number of inserts to generate</td>
<td><code>int</code></td>
<td>False</td>
<td><code>100</code></td>
<td></td>
</tr>
<tr>
<td><code>fields</code></td>
<td>Fields containing key name and function to run in json format</td>
<td><code>[]Field</code></td>
<td>False</td>
<td></td>
<td></td>
</tr>
</table></td>
</tr>
</table>


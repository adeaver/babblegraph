import React, { useState, useEffect } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import RadioGroup from '@material-ui/core/RadioGroup';

import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import { PrimaryRadio } from 'common/components/Radio/Radio';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { PrimaryButton, WarningButton } from 'common/components/Button/Button';
import { Heading3, Heading4, Heading6 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';
import BlogDisplay from 'common/components/BlogDisplay/BlogDisplay';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';

import {
    BlogPostMetadata,
    ContentNodeType,
    ContentNode,
    ContentNodeBody,
    Heading as HeadingContent,
    Paragraph as ParagraphContent,
    Image as ImageContent,
    Link as LinkContent,
    List as ListContent,
    ListType,
    getDefaultContentNodeForType,
} from 'common/api/blog/content';
import {
    getBlogContent,
    GetBlogContentResponse,
    updateBlogContent,
    UpdateBlogContentResponse,
} from 'AdminWeb/api/blog/blog';
import ImageUpload from 'AdminWeb/components/Blog/common/ImageUpload';

const styleClasses = makeStyles({
    blogContentEditorContainer: {
        padding: '10px',
    },
    editorRow: {
        minWidth: '100%',
        margin: '10px 0',
    },
    nodeButton: {
        margin: '0 5px',
    },
    nodeButtonVertical: {
        margin: '5px 0',
    }
});


type BlogContentEditorProps = {
    blogPostMetadata: BlogPostMetadata;
    urlPath: string;
}

const BlogContentEditor = (props: BlogContentEditorProps) => {
    const [ content, setContent ] = useState<Array<ContentNode>>([]);
    const [ isLoading, setIsLoading ] = useState<boolean>(true);
    const [ error, setError ] = useState<Error>(null);

    const handleRemoveContentAtIndex = (idx: number) => {
        setContent(content.filter((node: ContentNode, nodeIdx: number) => idx !== nodeIdx));
    }
    const handleAppendContentAtIndex = (newNode: ContentNode, idx: number) => {
        setContent(idx === content.length ? (
            content.concat(newNode)
        ) : (
            content.reduce((curr: Array<ContentNode>, nextNode: ContentNode, nodeIdx: number) => (
                nodeIdx === idx ? (
                    curr.concat([newNode, nextNode])
                ) : (
                    curr.concat(nextNode)
                )
            ), [])
        ));
    }
    const handleUpsertContentAtIndex = (newNode: ContentNode, idx: number) => {
        setContent(content.map((nextNode: ContentNode, nodeIdx: number) => (
            nodeIdx === idx ? newNode : nextNode
        )));
    }
    const handleSaveBlogContent = () => {
        setIsLoading(true);
        updateBlogContent({
            urlPath: props.urlPath,
            content: content,
        },
        (resp: UpdateBlogContentResponse) => {
            setIsLoading(false);
            if (!resp.success) {
                setError(new Error("Something went wrong"));
            }
        },
        (err: Error) => {
            setIsLoading(false);
            setError(err);
        });
    }

    useEffect(() => {
        getBlogContent({
            urlPath: props.urlPath,
        },
        (resp: GetBlogContentResponse) => {
            setIsLoading(false);
            setContent(resp.content || []);
        },
        (err: Error) => {
            setIsLoading(false);
            setError(err);
        });
    }, []);

    const classes = styleClasses();
    let body;
    if (isLoading) {
        body = <LoadingSpinner />;
    } else if (!!error) {
        body = (
            <Heading3 color={TypographyColor.Warning}>
                An error occurred.
            </Heading3>
        );
    } else {
        body = (
            <Grid container>
                <Grid className={classes.blogContentEditorContainer} item xs={6}>
                    <EditorComponent
                        urlPath={props.urlPath}
                        content={content}
                        handleSaveBlogContent={handleSaveBlogContent}
                        handleRemoveContentAtIndex={handleRemoveContentAtIndex}
                        handleAppendContentAtIndex={handleAppendContentAtIndex}
                        handleUpsertContentAtIndex={handleUpsertContentAtIndex} />
                </Grid>
                <Grid className={classes.blogContentEditorContainer} item xs={6}>
                    <DisplayCard>
                        <BlogDisplay
                            content={content}
                            metadata={props.blogPostMetadata} />
                    </DisplayCard>
                </Grid>
            </Grid>
        );
    }
    return (
        <div>
            {body}
        </div>
    );
}

type EditorComponentProps = {
    urlPath: string;
    content: ContentNode[];

    // TODO: add move content to index
    handleSaveBlogContent: () => void;
    handleRemoveContentAtIndex: (idx: number) => void;
    handleAppendContentAtIndex: (contentNode: ContentNode, idx: number) => void;
    handleUpsertContentAtIndex: (contentNode: ContentNode, idx: number) => void;
}

const EditorComponent = (props: EditorComponentProps) => {
    const [ newNodeType, setNewNodeType ] = useState<ContentNodeType>(ContentNodeType.Heading);

    const handleRadioFormChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setNewNodeType((event.target as HTMLInputElement).value as ContentNodeType);
    };

    const handleAddNewNode = () => {
        props.handleAppendContentAtIndex(getDefaultContentNodeForType(newNodeType), props.content.length);
    }
    const handleUpdateNode = (idx: number) => {
        return (newNode: ContentNode) => {
            props.handleUpsertContentAtIndex(newNode, idx);
        }
    }
    const handleDeleteNode = (idx: number) => {
        return () => {
            props.handleRemoveContentAtIndex(idx);
        }
    }

    const classes = styleClasses();
    return (
        <DisplayCard>
            {
                props.content.map((node: ContentNode, idx: number) => {
                    if (node.type === ContentNodeType.Heading) {
                        return (
                            <HeadingNodeEditor
                                body={node.body as HeadingContent}
                                handleDelete={handleDeleteNode(idx)}
                                handleUpdate={handleUpdateNode(idx)} />
                        );
                    } else if (node.type === ContentNodeType.Paragraph) {
                        return (
                            <ParagraphNodeEditor
                                body={node.body as ParagraphContent}
                                handleDelete={handleDeleteNode(idx)}
                                handleUpdate={handleUpdateNode(idx)} />
                        );
                    } else if (node.type === ContentNodeType.Image) {
                        return (
                            <ImageNodeEditor
                                urlPath={props.urlPath}
                                body={node.body as ImageContent}
                                handleDelete={handleDeleteNode(idx)}
                                handleUpdate={handleUpdateNode(idx)} />
                        )
                    } else if (node.type === ContentNodeType.Link) {
                        return (
                            <LinkNodeEditor
                                body={node.body as LinkContent}
                                handleDelete={handleDeleteNode(idx)}
                                handleUpdate={handleUpdateNode(idx)} />
                        )
                    } else if (node.type === ContentNodeType.List) {
                        return (
                            <ListNodeEditor
                                body={node.body as ListContent}
                                handleDelete={handleDeleteNode(idx)}
                                handleUpdate={handleUpdateNode(idx)} />
                        )
                    }
                })
            }
            <div className={classes.editorRow}>
                <Heading4 color={TypographyColor.Primary}>
                    Add a node </Heading4>
                <FormControl component="fieldset">
                    <RadioGroup aria-label="add-node-type" name="add-node-type1" value={newNodeType} onChange={handleRadioFormChange}>
                        <Grid container>
                            {
                                Object.values(ContentNodeType).map((option: ContentNodeType, idx: number) => (
                                    <EditorComponentAddNodeSelector
                                        key={`add-node-selector-${idx}`}
                                        value={option} />
                                ))
                            }
                        </Grid>
                    </RadioGroup>
                </FormControl>
            </div>
            <PrimaryButton
                type="submit"
                onClick={handleAddNewNode}>
                Submit
            </PrimaryButton>
            <PrimaryButton
                className={classes.nodeButton}
                type="submit"
                onClick={props.handleSaveBlogContent}>
                Save Updates
            </PrimaryButton>
        </DisplayCard>
    );
}

type EditorComponentAddNodeSelectorProps = {
    value: ContentNodeType | ListType;
}

const EditorComponentAddNodeSelector  = (props: EditorComponentAddNodeSelectorProps) => {
    return (
        <Grid item xs={6}>
            <FormControlLabel value={props.value} control={<PrimaryRadio />} label={props.value} />
        </Grid>
    )
}

type HeadingNodeEditorProps = {
    body: HeadingContent;
    handleUpdate: (node: ContentNode) => void;
    handleDelete: () => void;
}

const HeadingNodeEditor = (props: HeadingNodeEditorProps) => {
    const [ text, setText ] = useState<string>(props.body.text);
    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setText((event.target as HTMLInputElement).value);
    }

    const handleSubmit = () => {
        props.handleUpdate({
            type: ContentNodeType.Heading,
            body: {
                text: text,
            }
        })
    }

    const classes = styleClasses();
    return (
        <div>
            <Heading6 color={TypographyColor.Primary}>
                Edit Heading Node
            </Heading6>
            <PrimaryTextField
                className={classes.editorRow}
                value={text}
                label="Heading Content"
                onChange={handleChange}
                variant="outlined" />
            <PrimaryButton
                className={classes.nodeButton}
                type="submit"
                onClick={handleSubmit}>
                Submit
            </PrimaryButton>
            <WarningButton
                className={classes.nodeButton}
                type="submit"
                onClick={props.handleDelete}>
                Delete
            </WarningButton>
        </div>
    );
}

type ParagraphNodeEditorProps = {
    body: ParagraphContent;
    handleUpdate: (node: ContentNode) => void;
    handleDelete: () => void;
}

const ParagraphNodeEditor = (props: ParagraphNodeEditorProps) => {
    const [ text, setText ] = useState<string>(props.body.text);
    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setText((event.target as HTMLInputElement).value);
    }

    const handleSubmit = () => {
        props.handleUpdate({
            type: ContentNodeType.Paragraph,
            body: {
                text: text,
            }
        })
    }

    const classes = styleClasses();
    return (
        <div>
            <Heading6 color={TypographyColor.Primary}>
                Edit Paragraph Node
            </Heading6>
            <PrimaryTextField
                className={classes.editorRow}
                value={text}
                minRows={6}
                label="Paragraph Content"
                onChange={handleChange}
                variant="outlined"
                multiline />
            <PrimaryButton
                className={classes.nodeButton}
                type="submit"
                onClick={handleSubmit}>
                Submit
            </PrimaryButton>
            <WarningButton
                className={classes.nodeButton}
                type="submit"
                onClick={props.handleDelete}>
                Delete
            </WarningButton>
        </div>
    );
}

type ImageNodeEditorProps = {
    urlPath: string;
    body: ImageContent;
    handleUpdate: (node: ContentNode) => void;
    handleDelete: () => void;
}

const ImageNodeEditor = (props: ImageNodeEditorProps) => {
    const handleImageUpload = (i: ImageContent) => {
        props.handleUpdate({
            type: ContentNodeType.Image,
            body: i,
        })
    }

    const classes = styleClasses();
    return (
        <div>
            <ImageUpload
                label="Add an image"
                urlPath={props.urlPath}
                image={props.body}
                handleFileUpload={handleImageUpload} />
            <WarningButton
                className={classes.nodeButtonVertical}
                type="submit"
                onClick={props.handleDelete}>
                Delete
            </WarningButton>
        </div>
    );
}

type LinkNodeEditorProps = {
    body: LinkContent;
    handleUpdate: (node: ContentNode) => void;
    handleDelete: () => void;
}

const LinkNodeEditor = (props: LinkNodeEditorProps) => {
    const [ text, setText ] = useState<string>(props.body.text);
    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setText((event.target as HTMLInputElement).value);
    }

    const [ destinationURL, setDestinationURL ] = useState<string>(props.body.destinationUrl);
    const handleURLChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setDestinationURL((event.target as HTMLInputElement).value);
    }

    const handleSubmit = () => {
        props.handleUpdate({
            type: ContentNodeType.Link,
            body: {
                text: text,
                destinationUrl: destinationURL,
            }
        })
    }

    const classes = styleClasses();
    return (
        <div>
            <Heading6 color={TypographyColor.Primary}>
                Edit Link Node
            </Heading6>
            <PrimaryTextField
                className={classes.editorRow}
                value={text}
                minRows={6}
                label="Text"
                onChange={handleChange}
                variant="outlined"
                multiline />
            <PrimaryTextField
                className={classes.editorRow}
                value={destinationURL}
                label="Destination URL"
                onChange={handleURLChange}
                variant="outlined" />
            <PrimaryButton
                className={classes.nodeButton}
                type="submit"
                onClick={handleSubmit}>
                Submit
            </PrimaryButton>
            <WarningButton
                className={classes.nodeButton}
                type="submit"
                onClick={props.handleDelete}>
                Delete
            </WarningButton>
        </div>
    );
}

type ListNodeEditorProps = {
    body: ListContent;
    handleUpdate: (node: ContentNode) => void;
    handleDelete: () => void;
}

const ListNodeEditor = (props: ListNodeEditorProps) => {
    const [ items, setItems ] = useState<string[]>(props.body.items);
    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setItems((event.target as HTMLInputElement).value.split("\n"));
    }

    const [ listType, setListType ] = useState<ListType>(props.body.type);
    const handleRadioFormChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setListType((event.target as HTMLInputElement).value as ListType);
    };

    const handleSubmit = () => {
        props.handleUpdate({
            type: ContentNodeType.List,
            body: {
                items: items,
                type: listType,
            }
        })
    }

    const classes = styleClasses();
    return (
        <div>
            <Heading6 color={TypographyColor.Primary}>
                Edit List Node
            </Heading6>
            <PrimaryTextField
                className={classes.editorRow}
                value={items.join("\n")}
                minRows={6}
                label="Items (separated by newline)"
                onChange={handleChange}
                variant="outlined"
                multiline />
            <FormControl component="fieldset">
                <RadioGroup aria-label="add-node-type" name="add-node-type1" value={props.body.type} onChange={handleRadioFormChange}>
                    <Grid container>
                        {
                            Object.values(ListType).map((option: ListType, idx: number) => (
                                <EditorComponentAddNodeSelector
                                    key={`list-type-selector-${idx}`}
                                    value={option} />
                            ))
                        }
                    </Grid>
                </RadioGroup>
            </FormControl>
            <PrimaryButton
                className={classes.nodeButton}
                type="submit"
                onClick={handleSubmit}>
                Submit
            </PrimaryButton>
            <WarningButton
                className={classes.nodeButton}
                type="submit"
                onClick={props.handleDelete}>
                Delete
            </WarningButton>
        </div>
    );
}

export default BlogContentEditor;
